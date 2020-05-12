package taskrun

import (
	"context"
	"sync"
	"time"
)

type task struct {
	Name        string
	Index       int
	ChanIn      chan int
	List        *TaskList
	FnWork      func()
	WaitSeconds int
}

func (w *task) Start() {
	go func() {
		for {
			select {
			case value := <-w.ChanIn:
				if value == 0 {
					break
				}
				w.FnWork()
				w.List.Chan <- w.Index
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

type TaskList struct {
	lock  sync.Mutex
	Chan  chan int
	Works map[int]*task
}

func NewTaskList() *TaskList {
	return &TaskList{
		Chan:  make(chan int, 10),
		Works: make(map[int]*task),
	}
}
func (l *TaskList) AddTask(name string, waitSeconds int, fn func()) {

	l.lock.Lock()
	defer l.lock.Unlock()
	ch := make(chan int)
	index := len(l.Works)
	l.Works[index] = &task{
		List:        l,
		Index:       index,
		Name:        name,
		ChanIn:      ch,
		FnWork:      fn,
		WaitSeconds: waitSeconds,
	}

}
func (l *TaskList) Start(ctx context.Context) {
	go func() {
		waitStart := make(map[int]*Worker)

		for k, v := range l.Works {
			waitStart[k] = v
			v.Start()
		}
		i := 1

		for {
			i++

			select {
			case <-ctx.Done():
				for _, v := range waitStart {
					v.ChanIn <- v.Index
				}
				return
			case index := <-l.Chan:
				waitStart[index] = l.Works[index]
			default:
				for k, v := range waitStart {
					if (i % v.WaitSeconds) == 1 {
						if _, ok := waitStart[k]; ok {
							v.ChanIn <- 1
							delete(waitStart, k)
						}
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
