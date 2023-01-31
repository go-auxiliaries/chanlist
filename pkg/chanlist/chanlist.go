package chanlist

import (
	"github.com/go-auxiliaries/chanlist/pkg/caselist"
)

type List[T any] []chan T

var nilChannel = make(chan bool)

func New[T any](s int) *List[T] {
	out := make(List[T], s)
	return &out
}

func (e *List[T]) Len() int {
	return len(*e)
}

func (e *List[T]) Delete(channels ...chan T) int {
	l := *e
	cnt := 0
	for _, ch := range channels {
		for idx, el := range l {
			if el == ch {
				l[idx] = nil
				cnt += 1
				break
			}
		}
	}
	return cnt
}

func (e *List[T]) Init(chanLen int) *List[T] {
	l := *e
	for idx := range l {
		l[idx] = make(chan T, chanLen)
	}
	return e
}

func (e *List[T]) Write(idx int, item T) *List[T] {
	(*e)[idx] <- item
	return e
}

func (e *List[T]) Read(idx int) T {
	return <-(*e)[idx]
}

func (e *List[T]) Set(idx int, ch chan T) *List[T] {
	(*e)[idx] = ch
	return e
}

func (e *List[T]) Get(idx int) chan T {
	return (*e)[idx]
}

func (e *List[T]) Append(ch ...chan T) *List[T] {
	*e = append(*e, ch...)
	return e
}

func (e *List[T]) Clone() *List[T] {
	out := *e
	return &out
}

func (e *List[T]) AsList() *List[T] {
	return e
}

func (e *List[T]) CloneNoneEmpty() *List[T] {
	out := make(List[T], 0)
	for _, ch := range *e {
		if ch != nil {
			out = append(out, ch)
		}
	}
	return &out
}

func (e *List[T]) ToRecvCaseList() *caselist.List {
	out := caselist.New(0)
	for _, ch := range *e {
		if ch == nil {
			out.AppendRecv(nilChannel)
		} else {
			out.AppendRecv(ch)
		}
	}
	return out
}
