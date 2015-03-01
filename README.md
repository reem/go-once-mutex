# Once Mutex [![Build Status](https://travis-ci.org/reem/go-once-mutex.svg?branch=master)](https://travis-ci.org/reem/go-once-mutex)

> A Mutex offering one-time locking then infinite, concurrent, lock-free reads.

OnceMutex works like Mutex; it provides a synchronization primitive that can be
used to build other structures and manages other data in the absence of
generics. It is not meant for consumption by most users.

OnceMutex offers READ ONLY access after the initial lock - it is CRUCIAL that
the return value of `Lock` is ALWAYS checked and respected when mutation is
attempted, or data races could trivially arise.

This type is mostly useful for implementing lazily-evaluated types or other
similar primitives, since it can act as the implementation of blackholing.
OnceMutex is meant to be used within other libraries, and rarely by users
directly.

## Example

```go
package main

import oncem "github.com/reem/go-once-mutex"
import "fmt"

type Data struct {
    int x
}

func main() {
    // It is only legal to access data after thunk.Force has been called.
    data := &Data{0}
    once := oncem.NewOnceMutex()

    // "Expensive computation run!" will be printed once
    // some time after this.

    lockedbefore := once.Lock()

    if lockedbefore {
        panic("The once mutex was unexpectedly locked.")
    }

    go func() {
        data.x = 45
        once.Unlock()
    }()

    go func() {
        once.Lock() // Very cheap, just one atomic operation.
        fmt.Println("data.x:", data.x)
    }()
}
```

## Author

[Jonathan Reem](https://medium.com/@jreem) is the primary author and maintainer of future.

## License

MIT

