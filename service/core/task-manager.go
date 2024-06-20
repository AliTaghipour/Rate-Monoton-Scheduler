package core

import (
	"RateMonoticScheduler/model"
	"sync"
	"time"
)

type TaskManager struct {
	tasks         []*model.Task
	taskChan      chan *model.Task
	lock          sync.Mutex
	wg            sync.WaitGroup
	currentTaskId int
}

func NewTaskManager(tasks []*model.Task, taskChan chan *model.Task) *TaskManager {
	return &TaskManager{tasks: tasks, taskChan: taskChan, currentTaskId: 0}
}

func (t *TaskManager) Run() {
	t.wg.Add(len(t.tasks))
	for _, task := range t.tasks {
		go func(task *model.Task) {
			ticker := time.NewTicker(time.Duration(task.Period) * time.Second)
			for {
				select {
				case <-ticker.C:
					t.taskChan <- &model.Task{
						Id:       t.getNextTaskId(),
						Period:   task.Period,
						Duration: task.Duration,
					}
				default:

				}
			}
		}(task)
	}
	t.wg.Wait()
}

func (t *TaskManager) getNextTaskId() int {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.currentTaskId++
	return t.currentTaskId
}
