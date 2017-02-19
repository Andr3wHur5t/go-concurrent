package concurrent

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

/** High Level Interface **/

func TestBasicConcurrentQueueAfterAdd(t *testing.T) {
	concurrentTest := NewFunctionQueue(20) // Queue With 20 max queue entries
	assert.Equal(t, len(concurrentTest.workers), 0, "Should not start with workers")

	concurrentTest.SetWorkerCount(10)
	assert.Equal(t, len(concurrentTest.workers), 10, "Should be able to add workers")

	concurrentTest.SetWorkerCount(2)
	assert.Equal(t, len(concurrentTest.workers), 2, "Should be able to remove workers")

	concurrentTest.SetWorkerCount(125)
	assert.Equal(t, len(concurrentTest.workers), 125, "Should be able to add more workers")

	for _, o := range concurrentTest.workers {
		assert.Equal(t, o.isActive, true, "All workers should be active")
	}

	// We are putting a function on the queue to execute and opening channel to verify that execution completes.
	executionProof := make(chan bool)
	for i := 0; i <= 200; i++ {
		concurrentTest.queue <- func() {
			executionProof <- true
		}
		assert.Equal(t, <-executionProof, true, "Should execute function that was put on queue")
	}
}

func TestBasicConcurrentQueueAfterSubtract(t *testing.T) {
	concurrentTest := NewFunctionQueue(20) // Queue With 20 max queue entries
	assert.Equal(t, len(concurrentTest.workers), 0, "Should not start with workers")

	concurrentTest.SetWorkerCount(10)
	assert.Equal(t, len(concurrentTest.workers), 10, "Should be able to add workers")

	concurrentTest.SetWorkerCount(2)
	assert.Equal(t, len(concurrentTest.workers), 2, "Should be able to remove workers")

	concurrentTest.SetWorkerCount(125)
	assert.Equal(t, len(concurrentTest.workers), 125, "Should be able to add more workers")

	concurrentTest.SetWorkerCount(4)
	assert.Equal(t, len(concurrentTest.workers), 4, "Should be able to remove workers")

	for _, o := range concurrentTest.workers {
		assert.Equal(t, o.isActive, true, "All workers should be active")
	}

	// We are putting a function on the queue to execute and opening channel to verify that execution completes.
	executionProof := make(chan bool)
	for i := 0; i <= 200; i++ {
		concurrentTest.queue <- func() {
			executionProof <- true
		}
		assert.Equal(t, <-executionProof, true, "Should execute function that was put on queue")
	}
}

func _benchmarkExeCommon(workerCount int, b *testing.B) {
	concurrentTest := NewFunctionQueue(uint32(b.N) + 1)
	concurrentTest.SetWorkerCount(workerCount)

	var a float64
	loadFunc := func() {
		a = math.Pow(99999999, 99999999) / 100
	}

	didComplete := make(chan bool)
	compFunc := func() {
		didComplete <- true
	}

	for i := 0; i <= b.N; i++ {
		if i == b.N {
			concurrentTest.queue <- compFunc
		} else {
			concurrentTest.queue <- loadFunc
		}
	}

	<-didComplete
}

func BenchmarkQueueExecution1Workers(b *testing.B) {
	_benchmarkExeCommon(1, b)
}

func BenchmarkQueueExecution10Workers(b *testing.B) {
	_benchmarkExeCommon(10, b)
}

func BenchmarkQueueExecution100Workers(b *testing.B) {
	_benchmarkExeCommon(100, b)
}

func BenchmarkQueueExecution250Workers(b *testing.B) {
	_benchmarkExeCommon(250, b)
}

func BenchmarkQueueExecution500Workers(b *testing.B) {
	_benchmarkExeCommon(500, b)
}

func BenchmarkQueueExecution1000Workers(b *testing.B) {
	_benchmarkExeCommon(1000, b)
}

/** Worker Interface **/
