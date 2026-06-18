package system

import (
	"math/bits"

	"cache-protocols/cache"
	"cache-protocols/memory"
)

type System struct {
	Cache      cache.Cache
	Mem        memory.Memory
	AddrLayout AddressLayout
	Stats      Stats
}

type Stats struct {
	Hits uint64
	Misses uint64
	Reads uint64
	Writes uint64
}

type SplitAddress struct {
	AlignedAddress uint64
	Tag            uint64
	Index          uint64
	Offset         uint8
}

type AddressLayout struct {
	OffsetBits uint
	IndexBits  uint

	OffsetMask uint64
	IndexMask  uint64
}

func CacheRead(sys *System, address uint64, size uint64) []byte {
	sys.Stats.Reads++
	split := Decode(address, sys.AddrLayout)

	line, hit := sys.Cache.Read(split.Index, split.Tag)

	start := uint64(split.Offset)
	if hit {
		sys.Stats.Hits++
		return line[start : start+size]
	} else {
		sys.Stats.Misses++
	}

	block := sys.Mem.ReadBlock(split.AlignedAddress, sys.Cache.BlockSize)
	sys.Cache.Insert(split.Index, split.Tag, block)
	return block[start : start+size]
}

func Decode(address uint64, layout AddressLayout) SplitAddress {
	aligned := address & ^((uint64(1) << layout.OffsetBits) - 1)

	offset := (address & layout.OffsetMask)
	index := ((address >> layout.OffsetBits) & layout.IndexMask)
	tag := address >> uint64(layout.OffsetBits+layout.IndexBits)

	return SplitAddress{
		AlignedAddress: aligned,
		Tag:            tag,
		Index:          index,
		Offset:         uint8(offset),
	}
}

func GetAddressLayout(c cache.Cache) AddressLayout {
	if c.BlockSize == 0 || c.LineCount == 0 {
		panic("cache is improperly configured")
	}

	offsetBits := bits.TrailingZeros(uint(c.BlockSize))
	indexBits := bits.TrailingZeros(uint(c.LineCount))

	offsetMask := (1 << offsetBits) - 1
	indexMask := (1 << indexBits) - 1

	return AddressLayout{
		OffsetBits: uint(offsetBits),
		IndexBits:  uint(indexBits),
		OffsetMask: uint64(offsetMask),
		IndexMask:  uint64(indexMask),
	}
}
