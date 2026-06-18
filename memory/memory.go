package memory

import (
	"fmt"

	"cache-protocols/bitutils"
)

type Memory struct {
	Data []byte
	Size uint64
}

func NewMemory(size uint64) Memory {
	if !bitutils.IsPowerOfTwo(size) {
		panic("size of memory must be a power of two")
	}

	return Memory{
		Data: make([]byte, size),
		Size: size,
	}
}

func (m *Memory) ReadBlock(address uint64, size uint64) []byte {
	if (address + size) > m.Size {
		panic(fmt.Sprintf(
			"out of bounds: address=%d size=%d memSize=%d",
			address, size, m.Size,
		))
	}

	return m.Data[address : address+size]
}

func (m *Memory) Write(address uint64) {
}
