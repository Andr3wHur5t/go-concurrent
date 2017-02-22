package concurrent

type ConcurrentFunctionQueue struct {
	// Functions that are pushed here will be executed by the workers
	Queue FunctionQueue

	// Reference to our worker configuration
	Workers []*Worker
}

// Creates new Workers to meet the given count
func (q *ConcurrentFunctionQueue) spawnWorkers(workerCount uint32) {
	q.Workers = make([]*Worker, workerCount)

	// Create And Start Workers
	for i := (uint32)(0); i < workerCount; i++ {
		q.Workers[i] = NewWorker(q.Queue)
		go q.Workers[i].Start()
	}
}

// Starts our existing stopped workers
func (q *ConcurrentFunctionQueue) Start() {
	// Start Our Stopped Workers
	for _, worker := range q.Workers {
		if worker.isActive.Get() {
			continue
		}
		go worker.Start()
	}
}

// Stops All Workers
func (q *ConcurrentFunctionQueue) Stop() {
	for _, worker := range q.Workers {
		worker.Stop()
	}
}

// Creates A New Concurrent Function Queue with the max queue size as configured
func NewFunctionQueue(workerCount, maxQueueSize uint32) *ConcurrentFunctionQueue {
	q := new(ConcurrentFunctionQueue)
	q.Queue = make(FunctionQueue, maxQueueSize)
	q.spawnWorkers(workerCount)
	return q
}
