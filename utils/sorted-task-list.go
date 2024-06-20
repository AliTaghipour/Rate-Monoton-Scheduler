package utils

import (
	"RateMonoticScheduler/model"
	"errors"
	"slices"
	"sync"
)

type SortedTaskHandler interface {
	AddTask(task *model.Task) (result *model.Task)
	PopFirstTask() (*model.Task, error)
	GetTaskByID(id int) (*model.Task, error)
	ReplaceTask(task *model.Task, taskId int)
}

type SortedTaskHandlerImpl struct {
	tasks            []*model.Task
	freeProcessorIds []int
	lock             sync.Mutex
	size             int
}

func NewSortedTaskHandlerImpl(size int) SortedTaskHandler {
	var freeProcessorIds []int
	for i := 0; i < size; i++ {
		freeProcessorIds = append(freeProcessorIds, i+1)
	}
	return &SortedTaskHandlerImpl{size: size, freeProcessorIds: freeProcessorIds}
}

func (s *SortedTaskHandlerImpl) ReplaceTask(task *model.Task, taskId int) {
	for i, candidTask := range s.tasks {
		if candidTask.Id == taskId {
			task.ProcessorId = candidTask.ProcessorId
			s.tasks[i] = task
		}
	}
}

func (s *SortedTaskHandlerImpl) GetTaskByID(id int) (*model.Task, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var tasks []*model.Task
	var result *model.Task

	for _, task := range s.tasks {
		if task.Id != id {
			tasks = append(tasks, task)
		} else {
			result = task
		}
	}
	if result == nil {
		return nil, errors.New("task not found")
	}

	s.tasks = tasks
	s.freeProcessorIds = append(s.freeProcessorIds, result.ProcessorId)

	return result, nil
}

func (s *SortedTaskHandlerImpl) AddTask(task *model.Task) (result *model.Task) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.tasks) == s.size {
		if s.tasks[len(s.tasks)-1].Period < task.Period {
			result = task
			return
		}

		result = s.tasks[len(s.tasks)-1]
		task.ProcessorId = result.ProcessorId
		s.tasks = s.tasks[:len(s.tasks)-1]
	} else {
		task.ProcessorId = s.freeProcessorIds[0]
		s.freeProcessorIds = s.freeProcessorIds[1:]
	}

	s.tasks = append(s.tasks, task)
	slices.SortFunc(s.tasks, func(a, b *model.Task) int {
		if a.Period < b.Period {
			return -1
		} else if a.Period == b.Period {
			return 0
		} else {
			return 1
		}
	})
	return
}

func (s *SortedTaskHandlerImpl) PopFirstTask() (*model.Task, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.tasks) == 0 {
		return nil, errors.New("dklscn")
	}

	result := s.tasks[0]
	s.tasks = s.tasks[1:]

	s.freeProcessorIds = append(s.freeProcessorIds, result.ProcessorId)

	return result, nil

}
