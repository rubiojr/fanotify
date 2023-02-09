package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/s3rj1k/go-fanotify/fanotify"
	"golang.org/x/sys/unix"
)

var cmdArgs = os.Args[1:]
var watchPath string

func main() {
	flag.Parse()

	if watchPath == "" {
		fmt.Println("invalid path.")
		os.Exit(1)
	}

	if len(cmdArgs) == 0 {
		fmt.Println("invalid command.")
		os.Exit(1)
	}

	notify, err := fanotify.Initialize(
		unix.FAN_CLOEXEC|
			unix.FAN_CLASS_NOTIF|
			//unix.FAN_REPORT_FID| // doesn't work currently
			unix.FAN_UNLIMITED_QUEUE|
			unix.FAN_UNLIMITED_MARKS,
		os.O_RDONLY|
			unix.O_LARGEFILE|
			unix.O_CLOEXEC,
	)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	log.Printf("watching %s", watchPath)

	if err = notify.Mark(
		unix.FAN_MARK_ADD|
			unix.FAN_MARK_FILESYSTEM,
		unix.FAN_MODIFY|
			//unix.FAN_MOVE_SELF|
			//unix.FAN_ATTRIB|
			//unix.FAN_DELETE| doesn't work currently
			unix.FAN_CLOSE_WRITE,
		unix.AT_FDCWD,
		watchPath,
	); err != nil {
		log.Fatalf("%v\n", err)
	}

	f := func(notify *fanotify.NotifyFD) (string, error) {
		data, err := notify.GetEvent(os.Getpid())
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if data == nil {
			return "", nil
		}

		defer data.Close()

		path, err := data.GetPath()
		if err != nil {
			return "", err
		}

		dataFile := data.File()
		defer dataFile.Close()

		fInfo, err := dataFile.Stat()
		if err != nil {
			return "", err
		}

		mTime := fInfo.ModTime()

		if data.MatchMask(unix.FAN_CLOSE_WRITE) || data.MatchMask(unix.FAN_MODIFY) {
			if strings.HasPrefix(watchPath, path) {
				log.Printf("PID:%d %s - %v\n", data.GetPID(), path, mTime)
			}
			return path, nil
		}

		return "", fmt.Errorf("fanotify: unknown event")
	}

	for {
		path, err := f(notify)
		if err == nil && len(path) > 0 {
			if strings.HasPrefix(path, watchPath) {
				msg := SyncTask.WithArgs(context.Background())
				msg.Name = "sync"
				msg.OnceInPeriod(5 * time.Second)
				Queue.Add(msg)
			}
		}

		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}
}

func init() {
	flag.StringVar(&watchPath, "path", "", "Path to monitor")
}
