package concurrent

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

/** High Level Interface **/

func TestBasicConcurrentQueueAfterAdd(t *testing.T) {
	concurrentTest := NewFunctionQueue(45, 20) // Queue With 20 max queue entries, and 45 workers
	assert.Equal(t, len(concurrentTest.Workers), 45, "Should start with correct number of workers")

	for _, o := range concurrentTest.Workers {
		assert.Equal(t, o.isActive.Get(), true, "All workers should be active")
	}

	// We are putting a function on the queue to execute and opening channel to verify that execution completes.
	executionProof := make(chan bool)
	for i := 0; i <= 200; i++ {
		concurrentTest.Queue <- func() {
			executionProof <- true
		}
		assert.Equal(t, <-executionProof, true, "Should execute function that was put on queue")
	}
}

func TestBasicConcurrentQueueAfterSubtract(t *testing.T) {
	concurrentTest := NewFunctionQueue(4, 20) // Queue With 20 max queue entries, and 4 workers
	assert.Equal(t, len(concurrentTest.Workers), 4, "Should start with correct number of workers")

	for _, o := range concurrentTest.Workers {
		assert.Equal(t, o.isActive.Get(), true, "All workers should be active")
	}

	// We are putting a function on the queue to execute and opening channel to verify that execution completes.
	executionProof := make(chan bool)
	for i := 0; i <= 200; i++ {
		concurrentTest.Queue <- func() {
			executionProof <- true
		}
		assert.Equal(t, <-executionProof, true, "Should execute function that was put on queue")
	}
}

func _benchmarkExeCommon(workerCount uint32, b *testing.B) {
	// Spawn the number of workers to test, and make queue long enough to hold all bench entries
	concurrentTest := NewFunctionQueue(workerCount, uint32(b.N) + 1)

	var a float64
	loadFunc := func() {
		a = math.Pow(99999999, 99999999) / 100
	}

	didComplete := make(chan bool)
	compFunc := func() {
		didComplete <- true
	}

	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		if i == b.N {
			concurrentTest.Queue <- compFunc
		} else {
			concurrentTest.Queue <- loadFunc
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
