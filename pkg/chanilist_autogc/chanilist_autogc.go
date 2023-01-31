package chanilist_autogc

import (
	"sync"

	"github.com/go-auxiliaries/chanlist/pkg/chanilist"

	"github.com/go-auxiliaries/chanlist/pkg/caselist"
)

type GcCallBackFN[T any] func(removed []*T)

type ListInterface[T any, L any] interface {
	Len() int
	Write(idx int, item T)
	Read(idx int) T
	Set(idx int, ch chan T)
	AsList() []chan T
}

type List[T any, I chanilist.ChanListItem[T]] struct {
	chanList chanilist.List[T, I]
	lock     sync.RWMutex
	deleted  int
	gcLimit  int
}

func New[T any, I chanilist.ChanListItem[T]](s int) *List[T, I] {
	return &List[T, I]{
		chanList: *chanilist.New[T, I](s),
	}
}

func (e *List[T, I]) Len() int {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.chanList.Len()
}

func (e *List[T, I]) SetGCLimit(limit int) *List[T, I] {
	e.gcLimit = limit
	return e
}

func (e *List[T, I]) GetGCLimit() int {
	return e.gcLimit
}

func (e *List[T, I]) Write(idx int, item T) *List[T, I] {
	e.lock.RLock()
	defer e.lock.RUnlock()
	e.chanList.Write(idx, item)
	return e
}

func (e *List[T, I]) Read(idx int) T {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.chanList.Read(idx)
}

func (e *List[T, I]) Set(idx int, ch *I) *List[T, I] {
	e.lock.Lock()
	defer e.lock.Unlock()
	if ch == nil {
		e.incDeleted(1)
	}
	e.chanList.Set(idx, ch)
	return e
}

func (e *List[T, I]) Get(idx int) *I {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.chanList.Get(idx)
}

func (e *List[T, I]) runGarbageCollector() {
	e.lock.Lock()
	defer e.lock.Unlock()
	newChanList := *chanilist.New[T, I](0)
	for _, ch := range e.chanList {
		if ch != nil {
			newChanList = append(newChanList, ch)
		}
	}
	e.chanList = newChanList
	e.deleted = 0
}

func (e *List[T, I]) Append(items ...*I) *List[T, I] {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.chanList.Append(items...)
	return e
}

func (e *List[T, I]) Items() []*I {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.chanList
}

func (e *List[T, I]) Delete(items ...*I) int {
	e.lock.Lock()
	defer e.lock.Unlock()
	out := e.chanList.Delete(items...)
	e.incDeleted(out)
	return out
}

func (e *List[T, I]) incDeleted(n int) {
	e.deleted += n
	if e.deleted >= e.gcLimit {
		go e.runGarbageCollector()
	}
}

func (e *List[T, I]) ToRecvCaseList() *caselist.List {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.chanList.ToRecvCaseList()
}

func (e *List[T, I]) RLockBulk(fn func(l chanilist.List[T, I])) *List[T, I] {
	e.lock.RLock()
	defer e.lock.RUnlock()
	fn(e.chanList)
	return e
}
