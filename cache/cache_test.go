package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cache-protocols/cache"
)

func NewTestCache(size, blockSize uint64) cache.Cache {
	c := cache.Cache{
		Data:      make([]cache.CacheLine, size),
		LineCount: size,
		BlockSize: blockSize,
	}

	for i := range c.Data {
		c.Data[i] = cache.CacheLine{
			Valid: false,
			Tag:   0,
			Data:  make([]byte, blockSize),
		}
	}

	return c
}

func TestRead(t *testing.T) {
	t.Run("panics when index is out of bounds", func(t *testing.T) {
		c := NewTestCache(1, 1)

		assert.Panics(t, func() {
			_, _ = c.Read(2, 2)
		})
	})

	t.Run("returns a cache miss when the line is invalid", func(t *testing.T) {
		c := NewTestCache(1, 1)

		got, hit := c.Read(0, 0)

		assert.False(t, hit)
		assert.Nil(t, got)
	})

	t.Run("returns a cache miss when the tag does not match", func(t *testing.T) {
		c := NewTestCache(1, 1)
		c.Data[0].Valid = true

		got, hit := c.Read(0, 1)

		assert.False(t, hit)
		assert.Nil(t, got)
	})

	t.Run("returns a cache hit with correct data when matching tag and valid bit", func(t *testing.T) {
		c := NewTestCache(1, 1)
		c.Data[0].Valid = true

		got, hit := c.Read(0, 0)

		want := make([]byte, 1)
		assert.True(t, hit)
		assert.Equal(t, want, got)
	})
}

func TestInsert(t *testing.T) {
	t.Run("panics when the data and block size are mismatched", func(t *testing.T) {
		c := NewTestCache(1, 1)
		insertData := []byte("this sentence is too long for the cache")

		assert.Panics(t, func() {
			c.Insert(0, 0, insertData)
		})
	})

	t.Run("panics when the index is out of bounds", func(t *testing.T) {
		c := NewTestCache(1, 1)
		insertData := []byte("this sentence is too long for the cache")

		assert.Panics(t, func() {
			c.Insert(99, 0, insertData)
		})
	})

	t.Run("correctly inserts the data into the cache", func(t *testing.T) {
		c := NewTestCache(1, 1)
		insertData := []byte("w")

		c.Insert(0, 0, insertData)

		assert.Equal(t, insertData, c.Data[0].Data)
	})
}
