package handler

import (
	"testing"
	"time"

	"github.com/forgoes/logging"
	"github.com/forgoes/logging/diode"
	"github.com/forgoes/logging/writer"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRingBufferHandler(b *testing.B) {
	defer logging.Flush()
	defer logging.Close()

	w, err := writer.NewBufWriter(0, writer.NewStdoutWriter())
	assert.Nil(b, err)

	rh := NewRingBufferHandler(w, StdFormatter, 165536, diode.AlertFunc(func(int) {}), 100*time.Millisecond)

	logging.AddHandler(rh)
	defer logging.RemoveHandler(rh)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logging.Info().ALogf("buf stdout test: %d, %d, %d", 1, 2, 3)
		}
	})
}
