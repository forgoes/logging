package diode

import (
	"log"
	"sync/atomic"
	"unsafe"
)

// ManyToOne diode is optimal for many writers (go-routines B-n) and a single
// reader (go-routine A). It is not thread safe for multiple readers.
type ManyToOne struct {
	writeIndex uint64
	buffer     []unsafe.Pointer
	readIndex  uint64
	alerter    Alerter
}

// NewManyToOne creates a new diode (ring buffer). The ManyToOne diode
// is optimized for many writers (on go-routines B-n) and a single reader
// (on go-routine A). The alerter is invoked on the read's go-routine. It is
// called when it notices that the writer go-routine has passed it and wrote
// over data. A nil can be used to ignore alerts.
func NewManyToOne(size int, alerter Alerter) *ManyToOne {
	if alerter == nil {
		alerter = AlertFunc(func(int) {})
	}
	if size <= 0 {
		size = 1
	}

	d := &ManyToOne{
		buffer:  make([]unsafe.Pointer, size),
		alerter: alerter,
	}

	// Start write index at the value before 0
	// to allow the first write to use AddUint64
	// and still have a beginning index of 0
	d.writeIndex = ^d.writeIndex
	return d
}

// Set sets the data in the next slot of the ring buffer.
func (d *ManyToOne) Set(data GenericDataType) {
	for {
		writeIndex := atomic.AddUint64(&d.writeIndex, 1)
		idx := writeIndex % uint64(len(d.buffer))
		old := atomic.LoadPointer(&d.buffer[idx])

		if old != nil &&
			(*bucket)(old) != nil &&
			(*bucket)(old).seq > writeIndex-uint64(len(d.buffer)) {
			d.logCollision(writeIndex)
			continue
		}

		newBucket := &bucket{
			data: data,
			seq:  writeIndex,
		}

		if !atomic.CompareAndSwapPointer(&d.buffer[idx], old, unsafe.Pointer(newBucket)) {
			d.logCollision(writeIndex)
			continue
		}

		return
	}
}

func (d *ManyToOne) logCollision(writeIndex uint64) {
	// Sample collision logs to avoid flooding stderr in hot path.
	// Reuse the already-incremented writeIndex to avoid extra shared atomics.
	if writeIndex == 0 || writeIndex%1024 == 0 {
		log.Printf("Diode set collision (writeIndex=%d): consider using a larger diode", writeIndex)
	}
}

// TryNext will attempt to read from the next slot of the ring buffer.
// If there is no data available, it will return (nil, false).
func (d *ManyToOne) TryNext() (data GenericDataType, ok bool) {
	return tryNext(d.buffer, &d.readIndex, d.alerter)
}
