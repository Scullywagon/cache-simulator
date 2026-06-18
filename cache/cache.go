// Package cache does stuff
package cache

import (
	"cache-protocols/bitutils"
	"fmt"
)

type CacheLine struct {
	Valid bool
	Tag   uint64
	Data  []byte
}

type Cache struct {
	Data []CacheLine

	LineCount uint64
	BlockSize uint64
}

func NewCache(lineCount, blockSize uint64) Cache {
	if !bitutils.IsPowerOfTwo(lineCount) {
		panic("number of cache lines must be a power of two")
	}

	if !bitutils.IsPowerOfTwo(blockSize) {
		panic("block size must be a power of two")
	}

	lines := make([]CacheLine, lineCount)

	return Cache{
		Data:      lines,
		LineCount: lineCount,
		BlockSize: blockSize,
	}
}

func (cache Cache) Read(index uint64, tag uint64) ([]byte, bool) {
	if index >= cache.LineCount {
		panic(fmt.Sprintf("index was out of bound of cache size: %v", cache.LineCount))
	}

	line := cache.Data[index]

	if !line.Valid {
		return nil, false
	}
	if line.Tag != tag {
		return nil, false
	}

	return line.Data, true
}

func (cache *Cache) Insert(index uint64, tag uint64, blockData []byte) {
	if len(blockData) != int(cache.BlockSize) {
		panic("data and block size are mismatched")
	}

	if index >= cache.LineCount {
		panic(fmt.Sprintf("index was out of bound of cache size: %v", cache.LineCount))
	}

	cache.Data[index] = CacheLine{
		Valid: true,
		Tag:   tag,
		Data:  blockData,
	}
}
