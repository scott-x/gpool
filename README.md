# gpool
Library `gpool` implements a goroutine pool with fixed capacity, managing and recycling a massive number of goroutines

### Get Started

```bash
# install
go get github.com/scott-x/gpool
```

#### Example

calculate the num from 1 to 100

```go
package main

import (
	"fmt"
	"sync"

	"github.com/scott-x/gpool"
)

var sum = 0

type worker struct {
	mu sync.Mutex
}

func (w *worker) Do(i interface{}) {
	w.mu.Lock()
	sum += i.(int)
	w.mu.Unlock()
}

func main() {
	p := gpool.Init(10)

	w := worker{}
	//calculate the num from 1 to 100
	for i := 0; i <= 100; i++ {
		p.Do(&w, i)
	}
	p.Wait()
	fmt.Println(sum) //5050
}
```