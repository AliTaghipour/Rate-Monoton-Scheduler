package service

import "RateMonoticScheduler/service/core"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() {
	core.NewScheduler(1, 10)

}
