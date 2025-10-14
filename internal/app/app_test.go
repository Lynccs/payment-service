package app

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/stretchr/testify/require"
)

type fakeServer struct {
	startCh  chan struct{}
	startErr error
	stopped  bool
}

func (f *fakeServer) Start() error {
	close(f.startCh)
	if f.startErr != nil {
		return f.startErr
	}
	select {}
}

func (f *fakeServer) Stop(ctx context.Context) error {
	f.stopped = true
	return nil
}

func TestAppRun_ShutdownOnContextCancel(t *testing.T) {
	cfg := &config.Config{HTTPServer: config.HTTPServer{ShutdownTimeout: 2 * time.Second}}
	logger := slog.Default()
	fs := &fakeServer{startCh: make(chan struct{})}
	a := New(cfg, logger, fs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- a.Run(ctx) }()

	<-fs.startCh
	cancel()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("app.Run did not exit after context cancellation")
	}
	require.True(t, fs.stopped, "server Stop should be called")
}

func TestAppRun_ReturnsServerError(t *testing.T) {
	cfg := &config.Config{HTTPServer: config.HTTPServer{ShutdownTimeout: 2 * time.Second}}
	logger := slog.Default()
	expectedErr := errors.New("error starting server")
	fs := &fakeServer{startCh: make(chan struct{}), startErr: expectedErr}
	a := New(cfg, logger, fs)

	err := a.Run(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErr.Error())
}
