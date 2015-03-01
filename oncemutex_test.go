package oncemutex

import "testing"

type Data struct {
	x int
}

func TestMutateThenRead(t *testing.T) {
	once := NewOnceMutex()
	data := &Data{0}

	if once.Lock() {
		t.Fatal("OnceMutex reported it was previously locked on initialization.")
	}

	go func() {
		// We hold the lock and it was not locked already, so mutation is safe.
		data.x = 15

		once.Unlock() // Release the initial lock, so future locks are read-only.
	}()

	go func() {
		if !once.Lock() {
			t.Fatal("OnceMutex reported it was not previously locked after locking.")
		}

		if data.x != 15 {
			t.Fatal("OnceMutex reported it was not previously locked after locking.")
		}

		once.Unlock()
	}()
}
