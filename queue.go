package main

import (
	"log"
	"os/exec"

	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
)

// Create a queue.
var Queue = memqueue.NewQueue(&taskq.QueueOptions{
	Name:    "worker",
	Storage: taskq.NewLocalStorage(),
})

// Register a task.

var SyncTask = taskq.RegisterTask(&taskq.TaskOptions{
	Name: "sync",
	Handler: func() error {
		log.Print("changes detected, running sync")
		cmdPath, err := exec.LookPath(bin)
		if err != nil {
			return err
		}
		var cmd = exec.Command(cmdPath, cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			log.Printf("error running command: %s", err)
		}
		return nil
	},
})
