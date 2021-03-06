package oncemutex

import (
	"sync"
	"sync/atomic"
)

const (
	unused = 0
	locked = 1
	free   = 2
)

// A mutex which can only be locked once, but which provides
// very fast, lock-free, concurrent read-only locks after the
// first lock is over.
type OnceMutex struct {
	mu    sync.Mutex
	state uint32
}

func NewOnceMutex() *OnceMutex {
	return &OnceMutex{sync.Mutex{}, unused}
}

func (o *OnceMutex) Lock() (lockedbefore bool) {
	state := atomic.LoadUint32(&o.state)
	// The state is definitely free.
	if state == free {
		lockedbefore = true
		return
	}

	// The state is locked, or might have been unlocked already.
	if state == locked {
		// Once we have the lock check for a race.
		o.mu.Lock()

		// state could be free or could still be locked
		if atomic.LoadUint32(&o.state) == locked {
			// We acquired the mutex racily and incorrectly, unlock it
			// to allow the proper goroutine to acquire the lock.
			o.mu.Unlock()

			// Now when we acquire the lock the state will be free.
			o.mu.Lock()
		}

		o.mu.Unlock()

		lockedbefore = true
		return
	}

	if atomic.CompareAndSwapUint32(&o.state, unused, locked) {
		// Was unused, we are now locked.
		lockedbefore = false
		o.mu.Lock()
		return
	} else {
		// The previous state was changed from unused to locked.
		o.mu.Lock() // await the free state
		o.mu.Unlock()

		lockedbefore = true
		return
	}
}

func (o *OnceMutex) Unlock() {
	if atomic.CompareAndSwapUint32(&o.state, locked, free) {
		o.mu.Unlock()
		return
	}
}
