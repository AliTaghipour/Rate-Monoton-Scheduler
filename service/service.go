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
			Period:   15,
			Duration: 1,
		},
		{
			Period:   10,
			Duration: 0.5,
		},
		{
			Period:   20,
			Duration: 1,
		},
	}
	taskChan := core.NewScheduler(1, 10)
	core.NewTaskManager(tasks, taskChan).Run()

}
