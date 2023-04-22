package journalwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	// how many messages to buffer before sending
	BufferLen int
	// max time to wait before forwarding logs
	MaxBufferTime time.Duration
	// CloudWatch log group name you are forwarding into.
	// Kinda like the top folder in cloudwatch.
	LogGroupName string
	// CloudWatch individual log stream name. Often has the server and/or service
	// name in it.
	LogStreamName string
	// systemd unit name (the service name)
	Unit string
	// optional aws region, or will read from AWS_REGION
	Region string
}

func New(c *Config) *JournalWatch {
	if c.Region == "" {
		c.Region = os.Getenv("AWS_REGION")
	}
	return &JournalWatch{
		config:  c,
		stopped: true,
		wg:      sync.WaitGroup{},
	}
}

type JournalWatch struct {
	config  *Config
	stopped bool
	wg      sync.WaitGroup
	cw      *cloudwatchlogs.CloudWatchLogs
}

func (jw *JournalWatch) Start() error {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		log.Println("Error opening journal:", err)
		return err
	}

	defer journal.Close()

	err = journal.AddMatch("_SYSTEMD_UNIT=" + jw.config.Unit + ".service")
	if err != nil {
		log.Println("Error adding match:", err)
		return err
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(jw.config.Region),
	}))

	jw.cw = cloudwatchlogs.New(sess)

	logChan := make(chan *cloudwatchlogs.InputLogEvent)
	defer close(logChan)

	jw.stopped = false
	jw.wg.Add(2)
	go jw.putLogEvents(jw.config.LogGroupName, logChan)
	go jw.readJournal(journal, logChan)
	jw.wg.Wait()

	return nil
}

func (jw *JournalWatch) readJournal(journal *sdjournal.Journal, logChan chan<- *cloudwatchlogs.InputLogEvent) {
	for {
		n, err := journal.Next()
		if err != nil {
			log.Println("Error reading next entry:", err)
			continue
		}

		if n == 0 {
			journal.Wait(sdjournal.IndefiniteWait)
			continue
		}

		msg, err := journal.GetData("MESSAGE")
		if err != nil {
			log.Println("Error getting journal data:", err)
			continue
		}

		ts, err := journal.GetRealtimeUsec()
		if err != nil {
			log.Println("Error getting journald timestamp, skipping:", err)
			ts = uint64(time.Now().UnixNano())
		}

		t := time.Unix(0, int64(ts)*1000)

		event := &cloudwatchlogs.InputLogEvent{
			Message:   aws.String(msg),
			Timestamp: aws.Int64(t.UnixNano() / int64(time.Millisecond)),
		}

		logChan <- event
		if jw.stopped {
			jw.wg.Done()
			break
		}
	}
}

func (jw *JournalWatch) putLogEvents(logGroupName string, logChan <-chan *cloudwatchlogs.InputLogEvent) {
	bufferedEvents := []*cloudwatchlogs.InputLogEvent{}
	lastSent := time.Now()

	for event := range logChan {
		bufferedEvents = append(bufferedEvents, event)

		shouldSend := len(bufferedEvents) >= jw.config.BufferLen || time.Since(lastSent) >= jw.config.MaxBufferTime
		if !shouldSend {
			if jw.stopped {
				jw.wg.Done()
				break
			}
			continue
		}

		input := &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String("your_log_stream_name"),
			LogEvents:     bufferedEvents,
		}

		_, err := jw.cw.PutLogEvents(input)
		if err != nil {
			log.Println("Error sending log events to CloudWatch:", err)
			continue
		}

		bufferedEvents = []*cloudwatchlogs.InputLogEvent{}
		lastSent = time.Now()

		if jw.stopped {
			jw.wg.Done()
			break
		}
	}
}
