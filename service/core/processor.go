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
				fmt.Printf("processor [%d] - new task [%d]\n", s.id, task.Id)
				newTask = task
			default:

			}
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				s.lock.Lock()
				if s.currentTask != nil && s.currentTask.Duration > 0 {
					s.currentTask.Duration -= 0.001
					if s.currentTask.Duration <= 0 {
						fmt.Printf("processor [%d] - finished task [%d] \n", s.id, s.currentTask.Id)
						s.doneChannel <- DoneMessage{
							processorId: s.id,
							taskId:      s.currentTask.Id,
						}
					}
				}

				if newTask != nil && (s.currentTask == nil || s.currentTask.Id != newTask.Id) {
					previousId := 0
					if s.currentTask != nil {
						previousId = s.currentTask.Id
					}
					s.currentTask = newTask
					fmt.Printf("processor [%d] - new task with id [%d] took place of task [%d]\n", s.id, newTask.Id, previousId)

				}
				s.lock.Unlock()
				break

			default:

			}
		}
	}()
}
