package simulator_test

import (
	"math/rand/v2"
	"testing"

	"cache-protocols/simulator"
	"cache-protocols/system"

	"github.com/stretchr/testify/assert"
)

func NewStaticTestSystem(rng *rand.Rand, warmup simulator.WarmupType) system.System {
	cconf := simulator.CacheConfig{
		CacheLines: 10,
		BlockSize:  1,
	}
	sconf := simulator.SimConfig{
		WarmupType: warmup,
		RNG:        rng,
	}
	memSize := uint64(10)

	return simulator.NewSeededSystem(memSize, cconf, sconf)
}

func TestNewSeededSystem(t *testing.T) {
	t.Run("the system is filled deterministically when provided the same seed", func(t *testing.T) {
		rng1 := rand.New(rand.NewPCG(1, 0))
		rng2 := rand.New(rand.NewPCG(1, 0))

		got1 := NewStaticTestSystem(rng1, 0)
		got2 := NewStaticTestSystem(rng2, 0)

		assert.Equal(t, got1, got2)
	})

	t.Run("the system is filled deterministically when provided different seeds", func(t *testing.T) {
		rng1 := rand.New(rand.NewPCG(1, 0))
		rng2 := rand.New(rand.NewPCG(2, 0))

		got1 := NewStaticTestSystem(rng1, 0)
		got2 := NewStaticTestSystem(rng2, 0)

		assert.NotEqual(t, got1, got2)
	})

	t.Run("the memory is fully initialised", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 0))

		got := NewStaticTestSystem(rng, 0)

		assert.Len(t, got.Mem.Data, 10)
		for _, v := range got.Mem.Data {
			assert.NotEqual(t, byte(0), v)
		}
	})

	t.Run("cache is correctly configured", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 0))

		got := NewStaticTestSystem(rng, 0)

		assert.Equal(t, got.Cache.BlockSize, uint64(1))
		assert.Equal(t, got.Cache.LineCount, uint64(10))
		assert.Len(t, got.Cache.Data, 10)
	})

	t.Run("returns empty cache when cold warmup", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 0))

		got := NewStaticTestSystem(rng, 0)

		for _, v := range got.Cache.Data {
			assert.Equal(t, v.Data, []byte(nil))
			assert.Equal(t, v.Tag, uint64(0))
			assert.Equal(t, v.Valid, false)
		}
	})

	t.Run("random warmup populates cache", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 0))

		got := NewStaticTestSystem(rng, simulator.Random)

		validFound := false

		for _, line := range got.Cache.Data {
			if line.Valid {
				validFound = true
				break
			}
		}

		assert.True(t, validFound)
	})

	t.Run("partial warmup partially populates cache", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 0))

		got := NewStaticTestSystem(rng, simulator.Partial)

		validLines := 0

		for _, line := range got.Cache.Data {
			if line.Valid {
				validLines++
			}
		}

		assert.Greater(t, validLines, 0)
		assert.LessOrEqual(t, validLines, len(got.Cache.Data)/2)
	})

	t.Run("full warmup fully populates cache", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 0))

		got := NewStaticTestSystem(rng, simulator.Full)

		for i, line := range got.Cache.Data {
			assert.True(t, line.Valid)

			// a singular byte will convert into an int8
			assert.Equal(t, []byte{got.Mem.Data[i]}, line.Data)
		}
	})
}
