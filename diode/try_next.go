package diode

import (
	"sync/atomic"
	"unsafe"
)

func tryNext(buffer []unsafe.Pointer, readIndex *uint64, alerter Alerter) (GenericDataType, bool) {
	// Read a value from the ring buffer based on the readIndex.
	idx := *readIndex % uint64(len(buffer))
	result := (*bucket)(atomic.SwapPointer(&buffer[idx], nil))

	// When the result is nil that means the writer has not had the
	// opportunity to write a value into the diode. This value must be ignored
	// and the read head must not increment.
	if result == nil {
		return nil, false
	}

	// When the seq value is less than the current read index that means a
	// value was read from idx that was previously written but since has
	// been dropped. This value must be ignored and the read head must not
	// increment.
	if result.seq < *readIndex {
		return nil, false
	}

	// When the seq value is greater than the current read index that means a
	// value was read from idx that overwrote the value that was expected to
	// be at this idx. This happens when the writer has lapped the reader. The
	// reader needs to catch up to the writer, so it moves its write head to
	// the new seq, effectively dropping the messages that were not read in
	// between the two values.
	if result.seq > *readIndex {
		dropped := result.seq - *readIndex
		*readIndex = result.seq
		alerter.Alert(int(dropped))
	}

	// Only increment read index if a regular read occurred (where seq was
	// equal to readIndex) or a value was read that caused a fast-forward.
	*readIndex++
	return result.data, true
}
