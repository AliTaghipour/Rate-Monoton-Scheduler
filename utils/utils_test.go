package utils

import (
	"RateMonoticScheduler/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortedList(t *testing.T) {
	var (
		task1 = &model.Task{Id: 1, Period: 10}
		task2 = &model.Task{Id: 2, Period: 8}
		task3 = &model.Task{Id: 3, Period: 6}
		task4 = &model.Task{Id: 4, Period: 4}
	)

	taskHandler := NewSortedTaskHandlerImpl()

	taskHandler.AddTask(task2)
	taskHandler.AddTask(task3)

	task, err := taskHandler.PopFirstTask()
	assert.Nil(t, err)
	assert.Equal(t, task3, task)

	taskHandler.AddTask(task1)
	taskHandler.AddTask(task4)

	task, err = taskHandler.PopFirstTask()
	assert.Nil(t, err)
	assert.Equal(t, task4, task)

	task, err = taskHandler.PopFirstTask()
	assert.Nil(t, err)
	assert.Equal(t, task2, task)

	task, err = taskHandler.PopFirstTask()
	assert.Nil(t, err)
	assert.Equal(t, task1, task)

	task, err = taskHandler.PopFirstTask()
	assert.NotNil(t, err)

}
