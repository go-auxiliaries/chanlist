package caselist

import (
	"context"
	"reflect"
	"time"
)

type List struct {
	cases           []reflect.SelectCase
	defaultCase     int
	timeoutCase     int
	timeout         time.Duration
	timeoutBaseCtx  context.Context
	timeoutCancel   context.CancelFunc
	rechargeTimeout bool
}

func New(s int) *List {
	return &List{
		cases:       make([]reflect.SelectCase, s),
		defaultCase: -1,
		timeoutCase: -1,
	}
}

func (s *List) SetRechargeTimeout(val bool) *List {
	s.rechargeTimeout = val
	return s
}

func (s *List) GetRechargeTimeout() bool {
	return s.rechargeTimeout
}

func (s *List) SetRecv(idx int, ch any) *List {
	if idx == s.timeoutCase {
		s.resetTimeout()
	}
	if idx == s.defaultCase {
		s.defaultCase = -1
	}
	s.cases[idx] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	return s
}

func (s *List) SetSend(idx int, ch any, item any) *List {
	if idx == s.timeoutCase {
		s.resetTimeout()
	}
	if idx == s.defaultCase {
		s.defaultCase = -1
	}
	s.cases[idx] = reflect.SelectCase{Dir: reflect.SelectSend, Send: reflect.ValueOf(item), Chan: reflect.ValueOf(ch)}
	return s
}

func (s *List) Get(idx int) *reflect.SelectCase {
	return &s.cases[idx]
}

func (s *List) AppendRecv(ch any) *List {
	s.cases = append(s.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
	return s
}

func (s *List) AppendSend(ch any, item any) *List {
	s.cases = append(s.cases, reflect.SelectCase{Dir: reflect.SelectSend, Send: reflect.ValueOf(item), Chan: reflect.ValueOf(ch)})
	return s
}

func (s *List) AppendDefaultCase() *List {
	s.cases = append(s.cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	s.defaultCase = len(s.cases) - 1
	return s
}

func (s *List) Cancel() {
	if s.timeoutCancel != nil {
		s.timeoutCancel()
	}
}

func (s *List) AppendTimeoutCase(ctx context.Context, d time.Duration) *List {
	s.timeout = d
	s.timeoutBaseCtx = ctx
	ctx, s.timeoutCancel = context.WithTimeout(ctx, d)
	s.cases = append(s.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())})
	s.timeoutCase = len(s.cases) - 1
	return s
}

func (s *List) runSelect() (int, interface{}, bool) {
	s.rechargeTimeoutContext()
	idx, val, ok := reflect.Select(s.cases)
	if !ok {
		return idx, nil, false
	}
	return idx, val.Interface(), ok
}

func (s *List) Select() (int, interface{}, bool) {
	idx, val, ok := s.runSelect()
	if idx == s.defaultCase {
		return -1, val, ok
	}
	if idx == s.timeoutCase {
		return -1, val, ok
	}
	if s.defaultCase > idx {
		s.defaultCase -= 1
	}
	if s.timeoutCase > idx {
		s.timeoutCase -= 1
	}
	return idx, val, ok
}

func (s *List) Select2() (*reflect.SelectCase, interface{}, bool) {
	idx, val, ok := s.Select()
	if idx < 0 {
		return nil, val, ok
	}
	return &s.cases[idx], val, ok
}

func (s *List) SelectAutoCleanUp() (*reflect.SelectCase, interface{}, bool) {
	idx, val, ok := s.runSelect()
	if idx == s.defaultCase {
		return nil, val, ok
	}
	if idx == s.timeoutCase {
		return nil, val, ok
	}
	if s.defaultCase > idx {
		s.defaultCase -= 1
	}
	if s.timeoutCase > idx {
		s.timeoutCase -= 1
	}

	cs := s.cases[idx]

	last := len(s.cases) - 1
	for ; idx < last; idx += 1 {
		s.cases[idx] = s.cases[idx+1]
	}
	s.cases = s.cases[:len(s.cases)-1]
	return &cs, val, ok
}

func (s *List) Clone() *List {
	out := *s
	return &out
}

func (s *List) Len() int {
	return len(s.cases)
}

func (s *List) resetTimeout() {
	s.timeoutCase = -1
	s.timeout = time.Duration(0)
	s.timeoutCancel = nil
	s.timeoutBaseCtx = nil
}

func (s *List) rechargeTimeoutContext() {
	if !s.rechargeTimeout || s.timeoutCase < 0 {
		return
	}

	var ctx context.Context
	ctx, s.timeoutCancel = context.WithTimeout(s.timeoutBaseCtx, s.timeout)
	s.cases[s.timeoutCase].Chan = reflect.ValueOf(ctx.Done())
}
