package bitutils

func IsPowerOfTwo(v uint64) bool {
	return v != 0 && (v&(v-1)) == 0
}
