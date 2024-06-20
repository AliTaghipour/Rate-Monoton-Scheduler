package service

import (
	"RateMonoticScheduler/model"
	"RateMonoticScheduler/service/core"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() {
	tasks := []*model.Task{
		{
			Period:   5,
			Duration: 3,
		},
		{
			Period:   7,
			Duration: 0.9,
		},
		{
			Period:   11,
			Duration: 2,
		},
	}
	taskChan := core.NewScheduler(1, 10)
	core.NewTaskManager(tasks, taskChan).Run()

}
