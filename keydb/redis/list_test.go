package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *RedisTestSuite) TestInstancePush() {
	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("LPush", func(t *testing.T) {
		assert.NoError(t, s.instance.LPush(ctx, t.Name(), "val1"))
	})

	t.Run("RPush", func(t *testing.T) {
		assert.NoError(t, s.instance.RPush(ctx, t.Name(), "val2"))
	})
}

func (s *RedisTestSuite) TestInstancePop() {
	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("RPop", func(t *testing.T) {
		const expectedVal = "expectedVal"
		require.NoError(t, s.instance.RPush(ctx, t.Name(), expectedVal))
		var val string
		assert.NoError(t, s.instance.RPop(ctx, t.Name(), &val))
		assert.Equal(t, expectedVal, expectedVal)
	})

	t.Run("LPop", func(t *testing.T) {
		const expectedVal = "expectedVal"
		require.NoError(t, s.instance.RPush(ctx, t.Name(), expectedVal))
		var val string
		assert.NoError(t, s.instance.LPop(ctx, t.Name(), &val))
		assert.Equal(t, expectedVal, expectedVal)
	})
}

func (s *RedisTestSuite) TestInstance_RemoveFromList() {
	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("LRem", func(t *testing.T) {
		const expectedVal = "expectedVal"
		assert.NoError(t, s.instance.RemoveFromList(ctx, t.Name(), expectedVal))

		require.NoError(t, s.instance.RPush(ctx, t.Name(), expectedVal))
		assert.NoError(t, s.instance.RemoveFromList(ctx, t.Name(), expectedVal))
	})
}

func (s *RedisTestSuite) TestInstance_GetElementByPosition() {
	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("GetElementByPosition", func(t *testing.T) {
		var values = []int{1, 2, 3, 4, 5}
		for _, v := range values {
			require.NoError(t, s.instance.RPush(ctx, t.Name(), v))
		}

		var actualFirstPosVal int
		assert.NoErrorf(
			t, s.instance.GetElementByPosition(ctx, t.Name(), ListElementFirstPosition, &actualFirstPosVal),
			"error when get element with FIRST position",
		)
		assert.Equalf(t, values[0], actualFirstPosVal, "comparing of FIRST elements")

		var actualLastPosVal int
		assert.NoErrorf(
			t, s.instance.GetElementByPosition(ctx, t.Name(), ListElementLastPosition, &actualLastPosVal),
			"error when get element with LAST position",
		)
		assert.Equalf(t, values[len(values)-1], actualLastPosVal, "comparing of LAST elements")
	})
}

func (s *RedisTestSuite) TestInstance_GetList() {
	type (
		NestStruct struct {
			Str string
		}

		TestStruct struct {
			Int   int
			Str   string
			Float float32
			T     NestStruct
		}
	)

	var (
		ctx = context.Background()
		t   = s.T()
	)

	t.Run("int values", func(t *testing.T) {
		key := t.Name()
		expectedValues := []int{1, 2, 3}
		for _, val := range expectedValues {
			require.NoError(t, s.instance.RPush(ctx, key, val))
		}

		var values []int
		assert.NoError(t, s.instance.GetList(ctx, key, &values))
		assert.Equal(t, expectedValues, values)
	})

	t.Run("float values", func(t *testing.T) {
		key := t.Name()
		expectedValues := []float32{1.1, 2.2, 3.3}
		for _, val := range expectedValues {
			require.NoError(t, s.instance.RPush(ctx, key, val))
		}
		var values []float32
		assert.NoError(t, s.instance.GetList(ctx, key, &values))
		assert.Equal(t, expectedValues, values)
	})

	t.Run("string values", func(t *testing.T) {
		key := t.Name()
		expectedValues := []string{"val1", "val2", "val3"}
		for _, val := range expectedValues {
			require.NoError(t, s.instance.RPush(ctx, key, val))
		}
		var values []string
		assert.NoError(t, s.instance.GetList(ctx, key, &values))
		assert.Equal(t, expectedValues, values)
	})

	t.Run("struct values", func(t *testing.T) {
		key := t.Name()
		expectedValues := []TestStruct{
			{
				Int:   1,
				Str:   "1",
				Float: 1.1,
				T:     NestStruct{Str: "nested1"},
			},
			{
				Int:   2,
				Str:   "2",
				Float: 2.2,
				T:     NestStruct{Str: "nested2"},
			},
		}

		for _, val := range expectedValues {
			require.NoError(t, s.instance.RPush(ctx, key, val))
		}
		var values []TestStruct
		assert.NoError(t, s.instance.GetList(ctx, key, &values))
		assert.Equal(t, expectedValues, values)
	})
}
