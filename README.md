# gpool
Library `gpool` implements a goroutine pool with fixed capacity, managing and recycling a massive number of goroutines

### Get Started

```bash
# install
go get github.com/scott-x/gpool
```

#### Example1: calculate the num from 1 to 100

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

#### Example2: calculate the go files in $GOPATH

```go
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/scott-x/gpool"
)

var (
	sum               = 0
	dirPool, filePool *gpool.Pool
	wDir              dirWorker
	wFile             fileWorker
)

type dirWorker struct {
}

type fileWorker struct {
	mu sync.RWMutex
}

func (w *dirWorker) Do(i interface{}) {
	dir := i.(string)
	fls, _ := ioutil.ReadDir(dir)
	for _, v := range fls {
		name := v.Name()
		if name[0] == '.' || name[0] == '_' || name[0] == '$' || name[0] == '~' {
			continue
		}
		if v.IsDir() {
			dirPool.Do(&wDir, path.Join(dir, name))
		} else {
			filePool.Do(&wFile, path.Join(dir, name))
		}
	}
}

func (w *fileWorker) Do(i interface{}) {
	file := i.(string)
	ext := path.Ext(file)
	if ext == ".go" {
		w.mu.Lock()
		sum++
		w.mu.Unlock()
		// log.Println(file)
	}

}

func init() {
	dirPool = gpool.Init(10)
	filePool = gpool.Init(20)

	wDir = dirWorker{}
	wFile = fileWorker{}
}

func main() {
	t := time.Now()
	root := os.Getenv("GOPATH")
	fls, _ := ioutil.ReadDir(root)
	for _, v := range fls {
		name := v.Name()
		if name[0] == '.' || name[0] == '_' || name[0] == '$' || name[0] == '~' {
			continue
		}
		if v.IsDir() {
			dirPool.Do(&wDir, path.Join(root, name))
		} else {
			filePool.Do(&wFile, path.Join(root, name))
		}
	}
	filePool.Wait()
	dirPool.Wait()
	log.Println(time.Since(t))
	log.Printf("total golang files in $GOPATH: %d\n", sum)
}
```
Run 7 times, as you can see from the screenshot, it's very fast and accurate.

![](https://statics.scott-xiong.com/docusaurus/d38837ea66be405895981b77d9e26002.png)