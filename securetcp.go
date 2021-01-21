package fastsocks

import "sync"

const (
	bufSize = 1024
)

var bpool sync.Pool

func init() {
	bpool.New = func() interface{} {
		return make([]byte, bufSize)
	}
}

func bufferPoolGet() []byte {
	return bpool.Get().([]byte)
}

func bufferPoolPut(b []byte) {
	bpool.Put(b)
}
