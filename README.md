# Fanotify

Linux filesystem watcher using [fanotify](https://github.com/s3rj1k/go-fanotify)

## Installing

```
go install github.com/rubiojr/fanotify@latest
```

```
sudo setcap cap_sys_admin+eip ~/go/bin/fanotify
```

## Examples

Sync changes to `~/Documents` to a remote S3 bucket, using [rclone](https://rclone.org), when there are changes:

```
fanotify --path ~/Documents rclone sync $HOME/Documents s3:mybucket/Documents
```
