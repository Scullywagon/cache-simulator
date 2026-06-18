package system_test

import (
	"testing"

	"cache-protocols/cache"
	"cache-protocols/memory"
	"cache-protocols/system"

	"github.com/stretchr/testify/assert"
)

func MakeAddress(tag, index uint64, layout system.AddressLayout) uint64 {
	return (tag << (layout.IndexBits + layout.OffsetBits)) |
		(index << layout.OffsetBits)
}

func TestGetAddressLayout(t *testing.T) {
	t.Run("correctly determines address layout", func(t *testing.T) {
		cch := cache.Cache{
			Data:      []cache.CacheLine{},
			LineCount: 16,
			BlockSize: 32,
		}

		got := system.GetAddressLayout(cch)

		want := system.AddressLayout{
			OffsetBits: 5,
			IndexBits:  4,
			OffsetMask: 31,
			IndexMask:  15,
		}
		assert.Equal(t, want, got)
	})

	t.Run("panics when cache values are improperly set to zero", func(t *testing.T) {
		cch := cache.Cache{
			Data:      []cache.CacheLine{},
			LineCount: 0,
			BlockSize: 0,
		}

		assert.Panics(t, func() {
			system.GetAddressLayout(cch)
		})
	})
}

func TestCacheRead(t *testing.T) {
	t.Run("returns memory data and inserts into cache when cache miss", func(t *testing.T) {
		sys := &system.System{
			Cache: cache.Cache{
				Data: []cache.CacheLine{
					{Valid: false},
				},
				LineCount: 1,
				BlockSize: 1,
			},
			Mem: memory.Memory{
				Data: []byte{10},
				Size: 1,
			},
			AddrLayout: system.GetAddressLayout(cache.Cache{
				LineCount: 1,
				BlockSize: 1,
			}),
		}
		addr := MakeAddress(0, 0, sys.AddrLayout)

		got := system.CacheRead(sys, addr, 1)

		assert.Equal(t, []byte{10}, got)
		line, hit := sys.Cache.Read(0, 0)
		wantStats := system.Stats{
			Hits: 0,
			Misses: 1,
			Reads: 1,
			Writes: 0,
		}
		assert.Equal(t, wantStats, sys.Stats)
		assert.True(t, hit)
		assert.Equal(t, byte(10), line[0])
	})

	t.Run("returns cached data when cache hit", func(t *testing.T) {
		sys := &system.System{
			Cache: cache.Cache{
				Data: []cache.CacheLine{
					{
						Valid: true,
						Tag:   0,
						Data:  []byte{10},
					},
				},
				LineCount: 1,
				BlockSize: 1,
			},
			Mem: memory.Memory{}, // should not matter
			AddrLayout: system.AddressLayout{
				OffsetBits: 0,
				IndexBits:  0,
				OffsetMask: 0,
				IndexMask:  0,
			},
		}

		addr := MakeAddress(0, 0, sys.AddrLayout)

		got := system.CacheRead(sys, addr, 1)

		wantStats := system.Stats{
			Hits: 1,
			Misses: 0,
			Reads: 1,
			Writes: 0,
		}
		assert.Equal(t, wantStats, sys.Stats)
		assert.Equal(t, []byte{10}, got)
	})
}

func TestDecode(t *testing.T) {
	t.Run("correct bits and mask for regular values", func(t *testing.T) {
		layout := system.AddressLayout{
			OffsetBits: 2,
			IndexBits:  3,
			OffsetMask: 3,
			IndexMask:  7,
		}
		address := uint64(46)

		got := system.Decode(address, layout)

		want := system.SplitAddress{
			AlignedAddress: 44,
			Tag:            1,
			Index:          3,
			Offset:         2,
		}
		assert.Equal(t, want, got)
	})
	t.Run("correct bits and mask for minimum values", func(t *testing.T) {
		layout := system.AddressLayout{
			OffsetBits: 0,
			IndexBits:  1,
			OffsetMask: 0,
			IndexMask:  1,
		}
		address := uint64(46)

		got := system.Decode(address, layout)

		want := system.SplitAddress{
			AlignedAddress: 46,
			Tag:            23,
			Index:          0,
			Offset:         0,
		}
		assert.Equal(t, want, got)
	})
}
