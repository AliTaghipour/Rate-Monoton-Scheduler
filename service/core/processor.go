package core

import (
	"RateMonoticScheduler/model"
	"fmt"
	"sync"
	"time"
)

type DoneMessage struct {
	processorId int
	taskId      int
}

type Processor struct {
	id          int
	currentTask *model.Task
	lock        sync.Mutex
	taskChannel chan *model.Task
	doneChannel chan DoneMessage
}

func NewProcessor(id int, taskChannel chan *model.Task, doneChannel chan DoneMessage) *Processor {
	return &Processor{id: id, taskChannel: taskChannel, doneChannel: doneChannel}
}

func (s *Processor) process() {
	var newTask *model.Task

	go func() {
		for {
			select {
			case task := <-s.taskChannel:
				newTask = task
			default:

			}
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				s.lock.Lock()
				if s.currentTask != nil && s.currentTask.Duration > 0 {
					fmt.Printf("processor [%d] - processed task with id [%d]\n", s.id, s.currentTask.Id)
					s.currentTask.Duration -= 0.01
					if s.currentTask.Duration <= 0 {
						fmt.Printf("processor [%d] - finished task [%d] \n", s.id, s.currentTask.Id)
						s.doneChannel <- DoneMessage{
							processorId: s.id,
							taskId:      s.currentTask.Id,
						}
					}
				}

				if (s.currentTask == nil || s.currentTask.Duration <= 0) && newTask != nil && newTask.Duration > 0 {
					fmt.Printf("processor [%d] - new task with id [%d] took place of task [%d]\n", s.id, newTask.Id, s.currentTask.Id)
					s.currentTask = newTask

				}
				s.lock.Unlock()

			default:

			}
		}
	}()
}
