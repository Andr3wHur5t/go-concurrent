package concurrent

import "sync/atomic"

type AtomicBool int32

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32((*int32)(b)) == 1
}

func (b *AtomicBool) Set(newVal bool) {
	if newVal {
		atomic.StoreInt32((*int32)(b), 1)
	} else {
		atomic.StoreInt32((*int32)(b), 0)
	}
}
