package logging

import (
	"strconv"
	"sync"
)

var _bytesPool = &sync.Pool{
	New: func() interface{} {
		return &Bytes{
			bytes: make([]byte, 0, 1024),
		}
	},
}

func newBytes() *Bytes {
	b := _bytesPool.Get().(*Bytes)
	b.bytes = b.bytes[:0]
	return b
}

func putBytes(b *Bytes) {
	const max = 1 << 16
	if cap(b.bytes) > max {
		return
	}
	_bytesPool.Put(b)
}

type Bytes struct {
	bytes []byte
}

func (b *Bytes) AppendByte(s byte) {
	b.bytes = append(b.bytes, s)
}

func (b *Bytes) AppendInt(i int64) {
	b.bytes = strconv.AppendInt(b.bytes, i, 10)
}

func (b *Bytes) AppendString(s string) {
	b.bytes = append(b.bytes, s...)
}

func (b *Bytes) String() string {
	return string(b.bytes)
}
