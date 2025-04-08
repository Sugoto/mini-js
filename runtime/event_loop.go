package runtime

import (
	"sync"
	"time"
)

type Task struct {
	callback func()
	when     time.Time
}

type EventLoop struct {
	tasks []Task
	mu    sync.Mutex
}

func NewEventLoop() *EventLoop {
	el := &EventLoop{
		tasks: make([]Task, 0),
	}
	go el.run()
	return el
}

func (el *EventLoop) AddTask(callback func(), delay time.Duration) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.tasks = append(el.tasks, Task{
		callback: callback,
		when:     time.Now().Add(delay),
	})
}

func (el *EventLoop) Clear() {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.tasks = nil
}

func (el *EventLoop) run() {
	ticker := time.NewTicker(10 * time.Millisecond)
	for range ticker.C {
		el.mu.Lock()
		now := time.Now()
		remaining := make([]Task, 0)

		for _, task := range el.tasks {
			if now.After(task.when) {
				go task.callback()
			} else {
				remaining = append(remaining, task)
			}
		}

		el.tasks = remaining
		el.mu.Unlock()
	}
}
