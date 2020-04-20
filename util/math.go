package util

const (
	bitSize       = 32 << (^uint(0) >> 63)
	maxIntHeadBit = 1 << (bitSize - 2)
)

//ceil size  to power of 2,if size=15,return 16
func Ceil(size int) int {
	if size&maxIntHeadBit != 0 && size > maxIntHeadBit {
		panic("size is too large")
	}
	if size <= 2 {
		return 2
	}
	size--
	size = fillBits(size) //fill bits
	size++                //after fill,increment can make the size is the power of 2
	return size
}

//bit fill,promise that from n, all bits lower than the highest bit'1' of n are  all changed to bit'1',so named "fillBits"
func fillBits(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n
}
