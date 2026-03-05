package diode

import (
	"sync/atomic"
	"unsafe"
)

// GenericDataType is the data type the diodes operate on.
type GenericDataType unsafe.Pointer

// Alerter is used to report how many values were overwritten since the
// last write.
type Alerter interface {
	Alert(missed int)
}

// AlertFunc type is an adapter to allow the use of ordinary functions as
// Alert handlers.
type AlertFunc func(missed int)

// Alert calls f(missed)
func (f AlertFunc) Alert(missed int) {
	f(missed)
}

type bucket struct {
	data GenericDataType
	seq  uint64 // seq is the recorded write index at the time of writing
}

// OneToOne diode is meant to be used by a single reader and a single writer.
// It is not thread safe if used otherwise.
type OneToOne struct {
	buffer     []unsafe.Pointer
	writeIndex uint64
	readIndex  uint64
	alerter    Alerter
}

// NewOneToOne creates a new diode is meant to be used by a single reader and
// a single writer. The alerter is invoked on the read's go-routine. It is
// called when it notices that the writer go-routine has passed it and wrote
// over data. A nil can be used to ignore alerts.
func NewOneToOne(size int, alerter Alerter) *OneToOne {
	if alerter == nil {
		alerter = AlertFunc(func(int) {})
	}
	if size <= 0 {
		size = 1
	}

	return &OneToOne{
		buffer:  make([]unsafe.Pointer, size),
		alerter: alerter,
	}
}

// Set sets the data in the next slot of the ring buffer.
func (d *OneToOne) Set(data GenericDataType) {
	idx := d.writeIndex % uint64(len(d.buffer))

	newBucket := &bucket{
		data: data,
		seq:  d.writeIndex,
	}
	d.writeIndex++

	atomic.StorePointer(&d.buffer[idx], unsafe.Pointer(newBucket))
}

// TryNext will attempt to read from the next slot of the ring buffer.
// If there is no data available, it will return (nil, false).
func (d *OneToOne) TryNext() (data GenericDataType, ok bool) {
	return tryNext(d.buffer, &d.readIndex, d.alerter)
}
