package simulator

import (
	"math/rand/v2"

	"cache-protocols/cache"
	"cache-protocols/memory"
	"cache-protocols/system"
)

type WarmupType int

const (
	Cold WarmupType = iota
	Full
	Partial
	Random
)

type CacheConfig struct {
	CacheLines uint64
	BlockSize  uint64
}

type SimConfig struct {
	WarmupType WarmupType
	RNG        *rand.Rand
}

func NewSeededSystem(memorySize uint64, cacheConfig CacheConfig, simConfig SimConfig) system.System {
	c := cache.NewCache(cacheConfig.CacheLines, cacheConfig.BlockSize)
	m := memory.NewMemory(memorySize)

	sys := system.System{
		Cache: c,
		Mem:   m,
		Stats: system.Stats{},
	}

	sys.AddrLayout = system.GetAddressLayout(sys.Cache)

	fillMemory(&sys, simConfig.RNG)
	warmupCache(&sys, simConfig.RNG, simConfig.WarmupType)

	return sys
}

func fillMemory(sys *system.System, rng *rand.Rand) {
	for i := range sys.Mem.Data {
		sys.Mem.Data[i] = byte(rng.IntN(256))
	}
}

func warmupCache(sys *system.System, rng *rand.Rand, warmupType WarmupType) {
	switch warmupType {
	case Cold:
		warmupCold(sys)
	case Full:
		warmupFull(sys, rng)
	case Partial:
		warmupPartial(sys, rng)
	case Random:
		warmupRandom(sys, rng)
	}
}

func warmupCold(sys *system.System) {
	for i := range sys.Cache.Data {
		sys.Cache.Data[i].Valid = false
		sys.Cache.Data[i].Tag = 0
		sys.Cache.Data[i].Data = nil
	}
}

func warmupRandom(sys *system.System, rng *rand.Rand) {
	memoryBlocks := evalMemoryBlocks(sys)

	for range sys.Cache.Data {
		blockIndex := rng.IntN(int(memoryBlocks))
		addr := uint64(blockIndex) * sys.Cache.BlockSize

		blockData := sys.Mem.ReadBlock(addr, sys.Cache.BlockSize)

		split := system.Decode(addr, sys.AddrLayout)

		sys.Cache.Insert(split.Index, split.Tag, blockData)
	}
}

func warmupPartial(sys *system.System, rng *rand.Rand) {
	linesToFill := sys.Cache.LineCount / 2

	memoryBlocks := evalMemoryBlocks(sys)

	for range linesToFill {
		blockIndex := rng.IntN(int(memoryBlocks))
		addr := uint64(blockIndex) * sys.Cache.BlockSize

		blockData := sys.Mem.ReadBlock(addr, sys.Cache.BlockSize)

		split := system.Decode(addr, sys.AddrLayout)

		sys.Cache.Insert(split.Index, split.Tag, blockData)
	}
}

func warmupFull(sys *system.System, rng *rand.Rand) {
	memoryBlocks := evalMemoryBlocks(sys)

	for lineIndex := range sys.Cache.Data {
		var candidates []uint64

		for blockIndex := uint64(lineIndex); blockIndex < memoryBlocks; blockIndex += sys.Cache.LineCount {
			candidates = append(candidates, blockIndex)
		}

		chosen := candidates[rng.IntN(len(candidates))]
		addr := chosen * sys.Cache.BlockSize

		blockData := sys.Mem.ReadBlock(addr, sys.Cache.BlockSize)

		split := system.Decode(addr, sys.AddrLayout)

		sys.Cache.Insert(split.Index, split.Tag, blockData)
	}
}

func evalMemoryBlocks(sys *system.System) uint64 {
	return (sys.Mem.Size + sys.Cache.BlockSize - 1) / sys.Cache.BlockSize
}
