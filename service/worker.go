package service

type Worker struct {
	workerCount int
}

func NewWorker(workerCount int) *Worker {
	return &Worker{workerCount: workerCount}
}
