package logger

import (
	"bytes"
	"errors"
	"github.com/Lynccs/payment-service/internal/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	logger := SetupLogger("local")
	require.NotNil(t, logger)
}

func TestErr(t *testing.T) {
	err := errors.New("test err")
	attr := sl.Err(err)

	assert.Equal(t, "error", attr.Key)
	assert.Equal(t, err.Error(), attr.Value.String())
}

func TestLoggerOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	logger.Info("test info")

	if buf.Len() == 0 {
		t.Error("logger should write to output")
	}
}
