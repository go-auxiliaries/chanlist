package chanilist_autogc_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-auxiliaries/chanlist/pkg/chanilist_autogc"

	"github.com/stretchr/testify/assert"
)

type chanListItem struct {
	ch chan bool
}

func (i chanListItem) GetChannel() chan bool {
	return i.ch
}

func (i chanListItem) Send(val bool) {
	i.ch <- val
}

func (i chanListItem) Recv() bool {
	return <-i.ch
}

func Test_AutoGc_Delete(t *testing.T) {
	list := chanilist_autogc.New[bool, chanListItem](0)
	list.SetGCLimit(9)
	for x := 0; x < 10; x++ {
		list.Append(&chanListItem{make(chan bool, 1)})
	}

	assert.Equal(t, 10, list.Len())
	list.Delete(list.Items()[:9]...)
	time.Sleep(time.Second / 2)
	assert.Equal(t, 1, list.Len())
}

func Test_AutoGc_Set(t *testing.T) {
	list := chanilist_autogc.New[bool, chanListItem](0)
	list.SetGCLimit(9)
	for x := 0; x < 10; x++ {
		list.Append(&chanListItem{make(chan bool, 1)})
	}

	assert.Equal(t, 10, list.Len())
	for x := 0; x < 9; x++ {
		list.Set(x, nil)
	}
	time.Sleep(time.Second / 2)
	assert.Equal(t, 1, list.Len())
}

func Test_Recv_Select(t *testing.T) {
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	ch3 := make(chan bool, 1)

	list := chanilist_autogc.New[bool, chanListItem](0).Append(
		&chanListItem{ch1},
		&chanListItem{ch2},
		&chanListItem{ch3},
	)

	list.Write(0, false)
	idx, val, ok := list.ToRecvCaseList().Select()
	assert.Equal(t, 0, idx)
	assert.Equal(t, false, val)
	assert.Equal(t, true, ok)

	list.Write(1, true)
	idx, val, ok = list.ToRecvCaseList().Select()
	assert.Equal(t, 1, idx)
	assert.Equal(t, true, val)
	assert.Equal(t, true, ok)

	list.Write(2, true)
	idx, val, ok = list.ToRecvCaseList().AppendDefaultCase().Select()
	assert.Equal(t, 2, idx)
	assert.Equal(t, true, val)
	assert.Equal(t, true, ok)

	idx, val, ok = list.ToRecvCaseList().AppendDefaultCase().Select()
	assert.Equal(t, -1, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)

	now := time.Now().UTC()
	caseList := list.ToRecvCaseList().AppendTimeoutCase(context.Background(), time.Millisecond*100)
	defer caseList.Cancel()
	idx, val, ok = caseList.Select()
	assert.GreaterOrEqual(t, time.Now().UTC().Sub(now), time.Millisecond*100)
	assert.Equal(t, -1, idx)
	assert.Equal(t, nil, val)
	assert.Equal(t, false, ok)
}
