package chanilist

import (
	"github.com/go-auxiliaries/chanlist/pkg/caselist"
)

type ChanListItem[T any] interface {
	GetChannel() chan T
}

type List[T any, I ChanListItem[T]] []*I

var nilChannel = make(chan bool)

func New[T any, I ChanListItem[T]](s int) *List[T, I] {
	out := make(List[T, I], s)
	return &out
}

func (e *List[T, I]) Len() int {
	return len(*e)
}

func (e *List[T, I]) Delete(items ...*I) int {
	l := *e
	cnt := 0
	for _, ch := range items {
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

func (e *List[T, I]) Write(idx int, item T) *List[T, I] {
	(*((*e)[idx])).GetChannel() <- item
	return e
}

func (e *List[T, I]) Read(idx int) T {
	return <-(*((*e)[idx])).GetChannel()
}

func (e *List[T, I]) Set(idx int, item *I) *List[T, I] {
	(*e)[idx] = item
	return e
}

func (e *List[T, I]) Get(idx int) *I {
	return (*e)[idx]
}

func (e *List[T, I]) Append(items ...*I) *List[T, I] {
	*e = append(*e, items...)
	return e
}

func (e *List[T, I]) Clone() *List[T, I] {
	out := *e
	return &out
}

func (e *List[T, I]) AsList() *List[T, I] {
	return e
}

func (e *List[T, I]) CloneNoneEmpty() *List[T, I] {
	out := make(List[T, I], 0)
	for _, ch := range *e {
		if ch != nil {
			out = append(out, ch)
		}
	}
	return &out
}

func (e *List[T, I]) ToRecvCaseList() *caselist.List {
	out := caselist.New(0)
	for _, ch := range *e {
		if ch == nil {
			out.AppendRecv(nilChannel)
		} else {
			out.AppendRecv((*ch).GetChannel())
		}
	}
	return out
}
