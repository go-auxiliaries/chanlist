package chanlist_autogc

import (
	"sync"

	"github.com/go-auxiliaries/chanlist/pkg/caselist"
	"github.com/go-auxiliaries/chanlist/pkg/chanlist"
)

type GcCallBackFN func(removed []int)

type List[T any] struct {
	chanList   chanlist.List[T]
	lock       sync.RWMutex
	deleted    int
	gcLimit    int
	gcCallBack GcCallBackFN
}

func New[T any](s int) *List[T] {
	return &List[T]{
		chanList: *chanlist.New[T](s),
	}
}

func (e *List[T]) Len() int {
	return e.chanList.Len()
}

func (e *List[T]) SetGCCallback(cb GcCallBackFN) *List[T] {
	e.gcCallBack = cb
	return e
}

func (e *List[T]) SetGCLimit(limit int) *List[T] {
	e.gcLimit = limit
	return e
}

func (e *List[T]) GetGCLimit() int {
	return e.gcLimit
}

func (e *List[T]) Init(chanLen int) *List[T] {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.chanList.Init(chanLen)
	return e
}

func (e *List[T]) Write(idx int, item T) *List[T] {
	e.lock.RLock()
	defer e.lock.RUnlock()
	e.chanList.Write(idx, item)
	return e
}

func (e *List[T]) Read(idx int) T {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.chanList.Read(idx)
}

func (e *List[T]) Set(idx int, ch chan T) *List[T] {
	e.lock.Lock()
	defer e.lock.Unlock()
	if ch == nil {
		e.incDeleted(1)
	}
	e.chanList.Set(idx, ch)
	return e
}

func (e *List[T]) gcChanList() {
	e.lock.Lock()
	defer e.lock.Unlock()
	cb := e.gcCallBack
	if cb == nil {
		newChanList := make(chanlist.List[T], 0)
		for _, ch := range newChanList {
			if ch != nil {
				newChanList = append(newChanList, ch)
			}
		}
		e.deleted = 0
		return
	}
	newChanList := make(chanlist.List[T], 0)
	toRemove := make([]int, 0)
	for id, ch := range newChanList {
		if ch == nil {
			toRemove = append(toRemove, id)
		} else {
			newChanList = append(newChanList, ch)
		}
	}
	cb(toRemove)
	e.chanList = newChanList
	e.deleted = 0
}

func (e *List[T]) Append(ch ...chan T) *List[T] {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.chanList.Append(ch...)
	return e
}

func (e *List[T]) Delete(ch ...chan T) int {
	e.lock.Lock()
	defer e.lock.Unlock()
	out := e.chanList.Delete(ch...)
	e.incDeleted(out)
	return out
}

func (e *List[T]) incDeleted(n int) {
	e.deleted += n
	if e.deleted > e.gcLimit {
		go e.gcChanList()
	}
}

func (e *List[T]) ToRecvCaseList() *caselist.List {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.chanList.ToRecvCaseList()
}

func (e *List[T]) RLockBulk(fn func(l chanlist.List[T])) *List[T] {
	e.lock.RLock()
	defer e.lock.RUnlock()
	fn(e.chanList)
	return e
}
