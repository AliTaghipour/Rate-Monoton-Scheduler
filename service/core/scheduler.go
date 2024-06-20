package core

import (
	"RateMonoticScheduler/model"
	"RateMonoticScheduler/utils"
)

type Scheduler struct {
	workerCount   int
	channelSize   int
	taskChan      chan *model.Task
	waitingTasks  utils.SortedTaskHandler
	runningTasks  utils.SortedTaskHandler
	processorsMap map[int]*Processor

	doneChannel chan DoneMessage
}

func NewScheduler(workerCount int, channelSize int) chan *model.Task {
	channel := make(chan *model.Task, channelSize)
	doneChannel := make(chan DoneMessage)

	result := &Scheduler{
		workerCount:  workerCount,
		channelSize:  channelSize,
		taskChan:     channel,
		doneChannel:  doneChannel,
		waitingTasks: utils.NewSortedTaskHandlerImpl(channelSize),
		runningTasks: utils.NewSortedTaskHandlerImpl(workerCount),
	}

	go result.runScheduler()
	go result.runProcessors()

	return channel
}

func (s *Scheduler) runProcessors() {
	for i := 0; i < s.workerCount; i++ {
		c := make(chan *model.Task)
		processor := NewProcessor(i+1, c, s.doneChannel)
		s.processorsMap[i+1] = processor

		go processor.process()
	}
}

func (s *Scheduler) runScheduler() {
	for {
		select {
		case task := <-s.taskChan:
			s.processTask(task)
		case message := <-s.doneChannel:
			task, err := s.waitingTasks.PopFirstTask()
			if err == nil {
				s.runningTasks.ReplaceTask(task, message.taskId)
				s.runTask(task)
			}

		default:

		}
	}
}

func (s *Scheduler) processTask(task *model.Task) {
	result := s.runningTasks.AddTask(task)
	if result == task {
		s.waitingTasks.AddTask(task)
	} else {
		s.runTask(task)
	}
}

func (s *Scheduler) runTask(task *model.Task) {
	result := s.runningTasks.AddTask(task)
	if result != task {
		s.processorsMap[task.ProcessorId].taskChannel <- task
	}
}
