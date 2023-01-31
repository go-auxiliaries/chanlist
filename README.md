## Description ##

An API to ease select over dynamic channel list.
Under the hood it is using [reflect.Select](https://pkg.go.dev/reflect)

## Examples ##

### 1. Dynamic channel list ###

```go
package main

import (
	"context"
	"fmt"
	"github.com/go-auxiliaries/chanlist/pkg/chanlist"
	"time"
)

func main() {
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	ch3 := make(chan bool, 1)

	list := chanlist.New[bool](0).Append(ch1, ch2, ch3)

	ch4 := make(chan bool, 1)
	ch5 := make(chan bool, 1)

	ch1 <- true

	// select to receive data from the channels in the list
	idx, val, _ := list.ToRecvCaseList().Select()

	// channel at 0 got value 'true'
	fmt.Printf("channel at %d got value '%v'\n", idx, val)

	// add more channels to the list
	list = list.Append(ch4, ch5)

	ch4 <- false
	idx, val, _ = list.ToRecvCaseList().Select()

	// channel at 3 got value 'false'
	fmt.Printf("channel at %d got value '%v'\n", idx, val)

	idx, val, _ = list.ToRecvCaseList().AppendDefaultCase().Select()

	// channel at -1 got value '<nil>'
	fmt.Printf("channel at %d got value '%v'\n", idx, val)
}
```


### 2. Case list ###

TBD