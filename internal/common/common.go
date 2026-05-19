package common

import (
	"sync"
)

type Pixel struct {
	X, Y, Color byte
}

var slicePool = sync.Pool{
	New: func() any {
		s := make([]Pixel, 0, 1024)
		return &s
	},
}

func SetPixelPool() {
	PixelPool = slicePool.Get().(*[]Pixel)
}

func NilPixelPool() {
	PixelPool = nil
}

func PoolPut(pool *[]Pixel) {
	*pool = (*pool)[:0]
	slicePool.Put(pool)
}

var PixelPool *[]Pixel
