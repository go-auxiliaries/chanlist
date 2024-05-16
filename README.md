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
	"time"
	
    "github.com/go-auxiliaries/chanlist/pkg/chanlist"
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

	// None of channels have data, default case is going to fire
    idx, val, _ = list.ToRecvCaseList().AppendDefaultCase().Select()

	// channel at -1 got value '<nil>'
	fmt.Printf("channel at %d got value '%v'\n", idx, val)
}
```


### 2. Case list ###

```go
package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-auxiliaries/chanlist/pkg/caselist"

	"github.com/stretchr/testify/assert"
)

func main()  {
	boolCh := make(chan bool, 1)
	interfaceCh := make(chan interface{}, 1)

	caseList := caselist.New(0).
		AppendSend(boolCh, true).
		AppendSend(interfaceCh, true)

	boolCh <- true
	idx, val, ok := caseList.Select()
	if !ok || idx != 0 || val !=true {
		panic("it should read value from boolCh")
    }
}

```
