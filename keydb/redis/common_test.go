package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *RedisTestSuite) TestInstance_Common() {
	var (
		ctx = context.Background()
		t   = s.T()
	)

	t.Run("set", func(t *testing.T) {
		const val = "val"
		assert.NoError(t, s.instance.Set(ctx, t.Name(), val))
	})

	t.Run("get", func(t *testing.T) {
		const expectedVal = "val"
		require.NoError(t, s.instance.Set(ctx, t.Name(), expectedVal))

		var val string
		assert.NoError(t, s.instance.Get(ctx, t.Name(), &val))
		assert.Equal(t, expectedVal, val)
	})

	t.Run("get no value", func(t *testing.T) {
		var val string
		assert.ErrorIs(t, s.instance.Get(ctx, t.Name(), &val), ErrNoData)
	})

	t.Run("is_exists: not exists", func(t *testing.T) {
		ok, err := s.instance.IsExist(ctx, t.Name())
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("is_exists: exists", func(t *testing.T) {
		require.NoError(t, s.instance.Set(ctx, t.Name(), "str"))

		ok, err := s.instance.IsExist(ctx, t.Name())
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("del", func(t *testing.T) {
		const val = "val"
		require.NoError(t, s.instance.Set(ctx, t.Name(), val))
		assert.NoError(t, s.instance.Delete(ctx, t.Name()))
	})

	t.Run("del not existed key", func(t *testing.T) {
		assert.NoError(t, s.instance.Delete(ctx, t.Name()))
	})
}
