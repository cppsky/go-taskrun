package main

import (
	"context"
	"taskrun"
	"time"

	"github.com/cppsky/go-logger/logger"
)

func main() {
	l := taskrun.NewTaskList()
	l.AddTask("1", 3, func() {
		logger.Debug("run task 1")
		time.Sleep(5 * time.Second)
	})
	l.AddTask("2", 4, func() {
		logger.Debug("run task 2")
		time.Sleep(130 * time.Second)
	})
	l.AddTask("3", 5, func() {
		logger.Debug("run task 3")
		time.Sleep(40 * time.Second)
	})
	l.AddTask("4", 6, func() {
		logger.Debug("run task 4")
		time.Sleep(2 * time.Second)
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l.Start(ctx)
	time.Sleep(2 * time.Minute)
}
