# journalwatch

A minimal tool which forwards journald logs to AWS cloudwatch.

## Usage

```bash
journalwatch -h
```

```text
./journalwatch-x86
  -aws-region string
    	aws region name (default "us-west-2")
  -buffer-time int
    	log buffer max time in seconds to foward to cloudwatch regardless of buffer size (default 5)
  -h	print help
  -log-buffer int
    	log buffer max limit before forward to CloudWatch (default 10)
  -log-group string
    	(required) cloudwatch log group name
  -log-stream string
    	cloudwatch log stream name, defaults to hostname
  -unit string
    	(required) systemd unit name
```

Example systemd service

```unit file (systemd)

```
