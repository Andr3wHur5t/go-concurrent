package concurrent

type FunctionQueue chan func()

// Listens on a queue and executes functions
type Worker struct {
	// The queue the worker will pull tasks from
	queue       FunctionQueue

	// The flag that lets us kill the worker poll
	isActive    *AtomicBool

	// Signal to let the worker know it can pick up a new task
	didFinishTask chan bool
}

// Creates a worker which will execute tasks on the given function queue
func NewWorker(queue FunctionQueue) *Worker {
	w := new(Worker)
	w.queue = queue
	w.isActive = new(AtomicBool)
	w.didFinishTask = make(chan bool)
	return w
}

// Reports Task Completion
func (w *Worker) finishTask() {
	w.didFinishTask <- true
}

// Executes a task and reports when it's finished
func (w *Worker) startTask(task func()) {
	defer w.finishTask()
	task()
}

// Starts polling the configured queue and executes work
func (w *Worker) Start() {
	w.isActive.Set(true)
	for w.isActive.Get() {
		go w.startTask(<-w.queue)
		<-w.didFinishTask
	}
}

// Stops the worker from executing new tasks on the function queue
func (w *Worker) Stop() {
	w.isActive.Set(false)
}
