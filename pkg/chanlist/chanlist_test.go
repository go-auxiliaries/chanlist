package chanlist_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-auxiliaries/chanlist/pkg/chanlist"

	"github.com/stretchr/testify/assert"
)

func Test_ChanList_Recv_Select(t *testing.T) {
	list := chanlist.New[bool](3)
	list.Init(1)
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
