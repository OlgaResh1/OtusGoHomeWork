package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		const cacheSize = 3
		c := NewCache(cacheSize)
		var key Key
		var val string

		for i := 0; i < cacheSize; i++ {
			key = Key("Key_" + strconv.Itoa(i+1))
			val = "Value_" + strconv.Itoa(i+1)
			wasInCache := c.Set(key, val)
			require.False(t, wasInCache)

			readedVal, ok := c.Get(key)
			require.True(t, ok)
			require.Equal(t, val, readedVal)
		}
		// pop first element from cache
		wasInCache := c.Set(Key("Key_4"), "Value_4")
		require.False(t, wasInCache)

		_, ok := c.Get(Key("Key_1"))
		require.False(t, ok)

		// pop unused element
		readedVal, ok := c.Get(Key("Key_2"))
		require.True(t, ok)
		require.Equal(t, "Value_2", readedVal)

		wasInCache = c.Set(Key("Key_2"), "Value_2_new")
		require.True(t, wasInCache)

		wasInCache = c.Set(Key("Key_3"), "Value_3_new")
		require.True(t, wasInCache)

		readedVal, ok = c.Get(Key("Key_3"))
		require.True(t, ok)
		require.Equal(t, "Value_3_new", readedVal)

		wasInCache = c.Set(Key("Key_4_replaced"), "Value_4")
		require.False(t, wasInCache)

		_, ok = c.Get(Key("Key_4"))
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
