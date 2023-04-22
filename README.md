# `journalwatch`

A minimal tool which forwards journald logs to AWS cloudwatch.

The AWS cloudwatch forwarder is probably the right tool to use. Due to issues we've had related to logrotate, and
troubles naming the log group and log streams, this service was created.

`journalwatch` only listens on one systemd `-unit` / service at a time. Multiple instances of this app can be started
to cover multple log streams.

## Usage

Builds for linux x86 and arm64 are available in the Releases section of this repository.

The AWS CloudWatch log group and log stream will be created if they do not exist.


```bash
journalwatch -h
```

```text
./journalwatch-x86
  -aws-region string
    	aws region name (default "us-west-2")
  -buffer-time int
    	log buffer max time in seconds to foward to cloudwatch regardless of buffer size (default 8)
  -h	print help
  -log-buffer int
    	log buffer max limit before forward to CloudWatch (default 20)
  -log-group string
    	(required) cloudwatch log group name
  -log-stream string
    	cloudwatch log stream name, defaults to hostname
  -unit string
    	(required) systemd unit name
```

### How to run `journalwatch` using systemd in linux

Put a file at `/etc/systemd/system/journalwatch.service`.

(replace `/path/to/journalwatch` with the location of the `journalwatch` executable)

```unit file (systemd)
[Unit]
Description=Journalwatch Log Forwarder
After=network.target

[Service]
Type=simple
User=nobody
Group=nobody
ExecStart=/path/to/journalwatch
Restart=on-failure
RestartSec=3s

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl start journalwatch
sudo systemctl enable journalwatch
```

check for any issues:

```bash
sudo systemctl status journalwatch
sudo journalctl -u journalwatch -f
```
