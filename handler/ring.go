package handler

import (
	"context"
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
	return nil
}

func (r *RingBufferHandler) Close() error {
	r.cancel()
	<-r.done

	return nil
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

		if _, err := r.w.Write(b); err != nil {
			// TODO err handle
			e.Put()
			continue
		}

		if _, err := r.w.WriteString("\n"); err != nil {
			// TODO err handle
			e.Put()
			continue
		}

		e.Put()
	}
}
