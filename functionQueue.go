package concurrent

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
	action  func()
	didExit chan bool
}

// Creates a task for the given function
func NewTask(action func()) *Task {
	lt := new(Task)
	lt.action = action
	lt.didExit = make(chan bool)
	return lt
}

// Executes the task and reports the result
func (t *Task) Run() {
	defer func() {
		t.didExit <- true
	}()
	t.action()
}

/*** Worker ***/

// Listens on a queue and executes tasks
type Worker struct {
	queue       FunctionQueue
	isActive    bool
	currentTask *Task
}

// Creates a worker which will execute tasks on the given function queue
func NewWorker(queue FunctionQueue) *Worker {
	w := new(Worker)
	w.isActive = true
	w.queue = queue
	return w
}

// Stops the worker from executing new tasks on the function queue
func (w *Worker) Stop() {
	w.isActive = false
}

// Starts the worker on the function queue
func (w *Worker) Start() {
	w.isActive = true
	// Poll for work on our own go routine
	go func() {
		for w.isActive {
			// Pickup Work From the Queue
			w.currentTask = NewTask(<-w.queue)
			go w.currentTask.Run()
			// When We Report Exit Finish The Cycle and clean up
			<-w.currentTask.didExit
			w.currentTask = nil
		}
	}()
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
		q.workers[i].Start()
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
