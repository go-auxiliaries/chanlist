package chanilist_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-auxiliaries/chanlist/pkg/chanilist"

	"github.com/stretchr/testify/assert"
)

type chanListItem struct {
	ch chan bool
}

func (i chanListItem) GetChannel() chan bool {
	return i.ch
}

func Test_ChanList_Recv_Select(t *testing.T) {
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)
	ch3 := make(chan bool, 1)

	list := chanilist.New[bool, chanListItem](0).Append(
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
