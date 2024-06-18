package utils

import (
	"RateMonoticScheduler/model"
	"errors"
	"slices"
	"sync"
)

type SortedTaskHandler interface {
	AddTask(task *model.Task)
	PopFirstTask() (*model.Task, error)
}

type SortedTaskHandlerImpl struct {
	tasks []*model.Task
	lock  sync.Mutex
}

func NewSortedTaskHandlerImpl() SortedTaskHandler {
	return &SortedTaskHandlerImpl{}
}

func (s *SortedTaskHandlerImpl) AddTask(task *model.Task) {
	s.lock.Lock()
	defer s.lock.Unlock()

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
}

func (s *SortedTaskHandlerImpl) PopFirstTask() (*model.Task, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.tasks) == 0 {
		return nil, errors.New("dklscn")
	}

	result := s.tasks[0]
	s.tasks = s.tasks[1:]

	return result, nil

}
