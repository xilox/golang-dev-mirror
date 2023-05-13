package listeners

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hr3lxphr6j/bililive-go/src/configs"
	"github.com/hr3lxphr6j/bililive-go/src/instance"
	"github.com/hr3lxphr6j/bililive-go/src/live"
	livemock "github.com/hr3lxphr6j/bililive-go/src/live/mock"
)

func TestManagerAddAndRemoveListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.WithValue(context.Background(), instance.Key, &instance.Instance{})
	m := NewManager(ctx)
	backup := newListener
	newListener = func(ctx context.Context, live live.Live) Listener {
		ln := NewMockListener(ctrl)
		ln.EXPECT().Start().Return(nil)
		ln.EXPECT().Close()
		return ln
	}
	defer func() { newListener = backup }()
	l := livemock.NewMockLive(ctrl)
	l.EXPECT().GetLiveId().Return(live.ID("test")).Times(3)
	assert.NoError(t, m.AddListener(context.Background(), l))
	assert.Equal(t, ErrListenerExist, m.AddListener(context.Background(), l))
	ln, err := m.GetListener(context.Background(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, ln)
	assert.True(t, m.HasListener(context.Background(), "test"))
	assert.NoError(t, m.RemoveListener(context.Background(), "test"))
	assert.Equal(t, ErrListenerNotExist, m.RemoveListener(context.Background(), "test"))
	_, err = m.GetListener(context.Background(), "test")
	assert.Equal(t, ErrListenerNotExist, err)
	assert.False(t, m.HasListener(context.Background(), "test"))
}

func TestManagerStartAndClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.WithValue(context.Background(), instance.Key, &instance.Instance{
		Config: &configs.Config{
			RPC: configs.RPC{Enable: true},
		},
	})
	backup := newListener
	newListener = func(ctx context.Context, live live.Live) Listener {
		ln := NewMockListener(ctrl)
		ln.EXPECT().Start().Return(nil)
		ln.EXPECT().Close()
		return ln
	}
	defer func() { newListener = backup }()
	m := NewManager(ctx)
	assert.NoError(t, m.Start(ctx))
	for i := 0; i < 3; i++ {
		l := livemock.NewMockLive(ctrl)
		id := live.ID(fmt.Sprintf("test_%d", i))
		l.EXPECT().GetLiveId().Return(id).AnyTimes()
		assert.NoError(t, m.AddListener(ctx, l))
	}
	m.Close(ctx)
}
