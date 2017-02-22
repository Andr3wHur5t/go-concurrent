package concurrent

// Can replace this with Math.max, having this lets us have no dependencies.
func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}


type ConcurrentFunctionQueue struct {
	// Functions that are pushed here will be executed by the workers
	queue FunctionQueue

	// Reference to our worker configuration
	workers []*Worker
}

// Removes Existing Workers to meet the given count
func (q *ConcurrentFunctionQueue) reduceWorkers(workerCount int) {
	for i := len(q.workers) - 1; i >= workerCount; i-- {
		if q.workers[i] != nil {
			q.workers[i].Stop()
			q.workers[i] = nil
		}
	}

	q.workers = q.workers[:workerCount]
}

// Creates New Workers to meet the given count
func (q *ConcurrentFunctionQueue) spawnWorkers(workerCount int) {
	originalWorkerCount := len(q.workers)
	// Resize Workers Reference
	oldWorkers := q.workers
	q.workers = make([]*Worker, workerCount)
	copy(q.workers, oldWorkers)

	// Create And Start Workers
	for i := max(0, originalWorkerCount-1); i < workerCount; i++ {
		q.workers[i] = NewWorker(q.queue)
		go q.workers[i].Start()
	}
}

// Adds or removes workers to meet the given count
func (q *ConcurrentFunctionQueue) SetWorkerCount(workerCount int) {
	if len(q.workers) == workerCount {
		return
	}
	if len(q.workers) > workerCount {
		q.reduceWorkers(workerCount)
	} else {
		q.spawnWorkers(workerCount)
	}
}

// Creates A New Concurrent Function Queue with the max queue size as configured
func NewFunctionQueue(maxQueueSize uint32) *ConcurrentFunctionQueue {
	q := new(ConcurrentFunctionQueue)
	q.queue = make(FunctionQueue, maxQueueSize)
	return q
}
