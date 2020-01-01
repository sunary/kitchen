package number

import (
	"unsafe"
)

const (
	Phi         = 0x9E3779B9
	MagicNumber = 0x5f3759df
	OneHalf     = float32(0.5)
	ThreeHalfs  = float32(1.5)
)

// PhiMix Randomly changes each bit of the seed by key
func PhiMix(k uint64, mask uint64) uint64 {
	h := k * Phi
	return (h ^ (h >> 16)) & mask
}

// FastInverseSqrt Fast inverse square root, return 1/math.Sqrt(x)
func FastInverseSqrt(x float32) float32 {
	x2 := x * OneHalf
	y := x
	i := *(*int32)(unsafe.Pointer(&y))
	i = MagicNumber - i>>1
	y = *(*float32)(unsafe.Pointer(&i))
	return y * (ThreeHalfs - (x2 * y * y))
}
