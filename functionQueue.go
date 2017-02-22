package concurrent

import "github.com/tevino/abool"

type FunctionQueue chan func()

/*** Util ***/

// Can replace this with Math.max, having this lets us have no dependencies.
func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}


/*** Worker Task ***/

// Maintains Rudimentary State About A Workers Task
type Task struct {
	didExit chan bool
}

// Creates a task for the given function
func NewTask() *Task {
	lt := new(Task)
	lt.didExit = make(chan bool)
	return lt
}

// Executes the task and reports the result
func (t *Task) Run(action func()) {
	defer func() {
		t.didExit <- true
	}()
	action()
}

/*** Worker ***/

// Listens on a queue and executes tasks
type Worker struct {
	queue       FunctionQueue
	isActive    *abool.AtomicBool
	currentTask *Task
}

// Creates a worker which will execute tasks on the given function queue
func NewWorker(queue FunctionQueue) *Worker {
	w := new(Worker)
	w.isActive = abool.New()
	w.currentTask = NewTask()
	w.queue = queue
	return w
}
// Stops the worker from executing new tasks on the function queue
func (w *Worker) Stop() {
	w.isActive.SetTo(false)
}

// Starts the worker on the function queue
func (w *Worker) Start() {
	w.isActive.SetTo(true)

	// Poll for work on our own go routine
	for w.isActive.IsSet() {
		// Pickup Work From the Queue
		go w.currentTask.Run(<-w.queue)
		// When We Report Exit Finish The Cycle and clean up
		<-w.currentTask.didExit
	}
}

/*** Concurrent Function Queue ***/

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
