package handler

import (
	"context"
	"errors"
	"time"

	"github.com/forgoes/logging"
	"github.com/forgoes/logging/diode"
	"github.com/forgoes/logging/writer"
)

type Next func(*logging.Event) (bool, error)
type Handle func(*logging.Event) error

type Puller interface {
	diode.Diode

	Next() diode.GenericDataType
}

type RingBufferHandler struct {
	diode  diode.Diode
	puller Puller
	cancel context.CancelFunc
	done   chan struct{}
	format Formatter
	w      *writer.BufWriter
}

func NewRingBufferHandler(
	w *writer.BufWriter, f Formatter, size int, alert diode.AlertFunc, pullInterval time.Duration,
) *RingBufferHandler {
	ctx, cancel := context.WithCancel(context.Background())

	r := &RingBufferHandler{
		cancel: cancel,
		done:   make(chan struct{}),
		w:      w,
		format: f,
	}

	d := diode.NewManyToOne(size, alert)

	if pullInterval > 0 {
		r.puller = diode.NewPoller(
			d,
			diode.WithPollingInterval(pullInterval),
			diode.WithPollingContext(ctx),
		)
	} else {
		r.puller = diode.NewWaiter(
			d,
			diode.WithWaiterContext(ctx),
		)
	}

	go r.pull()

	return r
}

func (r *RingBufferHandler) Handle(e *logging.Event) (bool, error) {
	r.puller.Set(diode.GenericDataType(e.Clone()))
	return true, nil
}

func (r *RingBufferHandler) Flush() error {
	return r.w.Flush()
}

func (r *RingBufferHandler) Close() error {
	var closeErr error
	if err := r.w.Flush(); err != nil {
		closeErr = errors.Join(closeErr, err)
	}

	r.cancel()
	<-r.done

	if err := r.w.Close(); err != nil {
		closeErr = errors.Join(closeErr, err)
	}

	return closeErr
}

func (r *RingBufferHandler) pull() {
	defer close(r.done)

	for {
		d := r.puller.Next()
		if d == nil {
			return
		}

		e := (*logging.Event)(d)

		b, err := r.format(e)
		if err != nil {
			// TODO err handle
			e.Put()
			continue
		}

		// Write one complete line at once to avoid splitting a log line across
		// files when rotation happens between separate writes.
		//
		// Prefer appending in-place when capacity allows; fallback to copy.
		line := b
		if cap(line) > len(line) {
			line = line[:len(line):cap(line)]
			line = append(line, '\n')
		} else {
			tmp := make([]byte, len(b)+1)
			copy(tmp, b)
			tmp[len(b)] = '\n'
			line = tmp
		}
		if _, err := r.w.Write(line); err != nil {
			// TODO err handle
			e.Put()
			continue
		}

		e.Put()
	}
}
