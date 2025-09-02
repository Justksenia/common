package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *RedisTestSuite) TestTx() {
	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("committed transaction", func(t *testing.T) {
		var (
			val       = "val"
			channelID = "1"
		)

		tx, err := s.instance.Begin(ctx)
		require.NoError(t, err)

		require.NoError(t, tx.Set(ctx, channelID, val))
		require.NoError(t, tx.RPush(ctx, t.Name(), val))
		require.NoError(t, tx.Commit(ctx))

		var actualVal string
		assert.NoError(t, s.instance.Get(ctx, channelID, &actualVal))
		assert.Equal(t, val, actualVal)

		var actualListVal []string
		assert.NoError(t, s.instance.GetList(ctx, t.Name(), &actualListVal))
		assert.Equal(t, []string{val}, actualListVal)
	})

	t.Run("discarded transaction", func(t *testing.T) {
		var (
			val       = "val"
			channelID = "2"
		)

		tx, err := s.instance.Begin(ctx)
		require.NoError(t, err)

		require.NoError(t, tx.Set(ctx, channelID, val))
		require.NoError(t, tx.RPush(ctx, t.Name(), val))
		require.NoError(t, tx.Rollback(ctx))

		var actualVal string
		assert.ErrorIs(t, s.instance.Get(ctx, channelID, &actualVal), ErrNoData)
		assert.Equal(t, "", actualVal)

		var actualListVal []string
		assert.ErrorIs(t, s.instance.GetList(ctx, t.Name(), &actualListVal), ErrNoData)
	})
}
