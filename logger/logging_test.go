package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLoggingSingleton(t *testing.T) {
	const (
		iterations = 10
	)
	addresses := make([]*zap.Logger, iterations)
	firstAddress := make([]*zap.Logger, iterations)
	firstLogger, err := New()
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		firstAddress[i] = firstLogger
		var l *zap.Logger
		l, err = New()
		require.NoError(t, err)
		addresses[i] = l
	}

	assert.ElementsMatch(t, firstAddress, addresses)
}
