package caselist_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-auxiliaries/chanlist/pkg/caselist"

	"github.com/stretchr/testify/assert"
)

func Test_Clone(t *testing.T) {
	boolCh := make(chan bool)
	interfaceCh := make(chan interface{})

	caseList := caselist.New(0).
		AppendSend(boolCh, true).
		AppendSend(interfaceCh, true)

	cloneList := caseList.Clone()
	caseList = caseList.AppendDefaultCase()
	assert.Equal(t, 2, cloneList.Len())
	assert.Equal(t, 3, caseList.Len())
	cloneList = cloneList.AppendDefaultCase()
	cloneList.SetRecv(0, make(chan int))
	assert.Equal(t, 3, cloneList.Len())
	assert.Equal(t, 3, caseList.Len())

	assert.NotEqual(t, caseList.Get(0), cloneList.Get(0))

}

func Test_Send_Select(t *testing.T) {
	boolCh := make(chan bool, 1)
	interfaceCh := make(chan interface{}, 1)

	caseList := caselist.New(0).
		AppendSend(boolCh, true).
		AppendSend(interfaceCh, true)

	boolCh <- true
	idx, val, ok := caseList.Select()
	assert.Equal(t, 1, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)

	<-boolCh

	idx, val, ok = caseList.Select()
	assert.Equal(t, 0, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)

	caseListWithDefault := caseList.Clone().AppendDefaultCase()

	<-boolCh

	idx, val, ok = caseListWithDefault.Select()
	assert.Equal(t, 0, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)

	idx, val, ok = caseListWithDefault.Select()
	assert.Equal(t, 2, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)

	now := time.Now().UTC()
	selectList := caseList.AppendTimeoutCase(context.Background(), time.Millisecond*100)
	defer caseList.Cancel()

	idx, val, ok = selectList.Select()
	assert.GreaterOrEqual(t, time.Now().UTC().Sub(now), time.Millisecond*100)
	assert.Equal(t, -1, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
}

func Test_Send_NoAutoRecharge(t *testing.T) {
	boolCh := make(chan bool, 1)
	interfaceCh := make(chan interface{}, 1)

	caseList := caselist.New(0).
		AppendSend(boolCh, true).
		AppendTimeoutCase(context.Background(), time.Millisecond*100).
		AppendSend(interfaceCh, false)

	assert.Equal(t, 3, caseList.Len())

	defer caseList.Cancel()

	now := time.Now().UTC()
	caseList.SelectAutoCleanUp()
	assert.Equal(t, 2, caseList.Len())
	caseList.SelectAutoCleanUp()
	assert.Equal(t, 1, caseList.Len())

	cs, val, ok := caseList.SelectAutoCleanUp()
	assert.Equal(t, (*reflect.SelectCase)(nil), cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 1, caseList.Len())
	assert.GreaterOrEqual(t, time.Now().UTC().Sub(now), time.Millisecond*100)

	now = time.Now().UTC()
	cs, val, ok = caseList.SelectAutoCleanUp()
	assert.Equal(t, (*reflect.SelectCase)(nil), cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 1, caseList.Len())
	assert.Less(t, time.Now().UTC().Sub(now), time.Millisecond*100)
}

func Test_Send_AutoRecharge(t *testing.T) {
	boolCh := make(chan bool, 1)
	interfaceCh := make(chan interface{}, 1)

	caseList := caselist.New(0).
		SetRechargeTimeout(true).
		AppendSend(boolCh, true).
		AppendTimeoutCase(context.Background(), time.Millisecond*100).
		AppendSend(interfaceCh, false)

	assert.Equal(t, 3, caseList.Len())

	defer caseList.Cancel()

	now := time.Now().UTC()
	caseList.SelectAutoCleanUp()
	assert.Equal(t, 2, caseList.Len())
	caseList.SelectAutoCleanUp()
	assert.Equal(t, 1, caseList.Len())

	cs, val, ok := caseList.SelectAutoCleanUp()
	assert.Equal(t, (*reflect.SelectCase)(nil), cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 1, caseList.Len())
	assert.GreaterOrEqual(t, time.Now().UTC().Sub(now), time.Millisecond*100)

	now = time.Now().UTC()
	cs, val, ok = caseList.SelectAutoCleanUp()
	assert.Equal(t, (*reflect.SelectCase)(nil), cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 1, caseList.Len())
	assert.GreaterOrEqual(t, time.Now().UTC().Sub(now), time.Millisecond*100)
}

func Test_Send_SelectAutoCleanUp(t *testing.T) {
	boolCh := make(chan bool, 1)
	interfaceCh := make(chan interface{}, 1)

	caseList := caselist.New(0).
		AppendSend(boolCh, true).
		AppendSend(interfaceCh, false).
		AppendDefaultCase()

	interfaceCh <- false
	assert.Equal(t, 3, caseList.Len())

	cs, val, ok := caseList.SelectAutoCleanUp()
	assert.Equal(t, &reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(boolCh),
		Send: reflect.ValueOf(true),
	}, cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 2, caseList.Len())

	<-interfaceCh

	cs, val, ok = caseList.SelectAutoCleanUp()
	assert.Equal(t, &reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(interfaceCh),
		Send: reflect.ValueOf(false),
	}, cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 1, caseList.Len())

	cs, val, ok = caseList.SelectAutoCleanUp()
	assert.Equal(t, (*reflect.SelectCase)(nil), cs)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
	assert.Equal(t, 1, caseList.Len())
}
