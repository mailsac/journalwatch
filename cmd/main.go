package main

import (
	"flag"
	"journalwatch/journalwatch"
	"log"
	"os"
	"time"
)

var unitFileName = flag.String("unit", "", "(required) systemd unit name")
var logGroup = flag.String("log-group", "", "(required) cloudwatch log group name")
var logStream = flag.String("log-stream", "", "cloudwatch log stream name, defaults to hostname")

var awsRegion = flag.String("aws-region", "us-west-2", "aws region name")

var logBufferLimit = flag.Int("log-buffer", 10, "log buffer max limit before forward to CloudWatch")
var timeoutSeconds = flag.Int("buffer-time", 5, "log buffer max time in seconds to foward to cloudwatch regardless of buffer size")
var help = flag.Bool("h", false, "print help")

func main() {
	flag.Parse()

	if *logStream == "" {
		hn, err := os.Hostname()
		if hn == "" {
			estr := err.Error()
			logStream = &estr
		} else {
			logStream = &hn
		}
	}
	if *logGroup == "" || *logStream == "" || *help {
		flag.PrintDefaults()
		return
	}
	c := &journalwatch.Config{
		BufferLen:     *logBufferLimit,
		MaxBufferTime: time.Second * time.Duration(*timeoutSeconds),
		LogGroupName:  *logGroup,
		LogStreamName: *logStream,
		Unit:          *unitFileName,
		Region:        *awsRegion,
	}
	jw := journalwatch.New(c)

	log.Println("starting journalwatch", *c)

	err := jw.Start()
	if err != nil {
		log.Println("journalwatch failed to initialize", err)
		os.Exit(1)
	}
	log.Println("journalwatch initialized successfully")
}
