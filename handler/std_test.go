package handler

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgoes/logging"
)

func TestStderrHandlerWithLineFormatter(t *testing.T) {
	defer logging.Flush()

	logging.SetLevel(logging.DEBUG)

	logging.EnableCaller(logging.INFO)
	defer logging.DisableCaller(logging.INFO)
	logging.EnableCaller(logging.ERROR)
	defer logging.DisableCaller(logging.ERROR)

	logging.EnableStack(logging.ERROR)
	defer logging.DisableStack(logging.ERROR)

	h := NewStderrHandler(StdFormatter)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	logging.Debug().Log()
	logging.Debug().Logf("stderr")

	logging.Info().Log()
	logging.Info().Logf("stderr")

	t.Run("parallel info", func(t *testing.T) {
		var w sync.WaitGroup
		for i := 0; i < 10; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()

				logging.Info().Logf("test parallel")
			}(i)
		}
		w.Wait()
	})

	logging.Warn().Log()
	logging.Warn().Logf("error")

	logging.Error().E(errors.New("with error")).Logf("error")

	logger := logging.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(logging.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(logging.ERROR)

	logger.Error().Logf("error")

	logger.EnableCaller(logging.PANIC)
	logger.EnableStack(logging.PANIC)
	assert.Panics(t, func() {
		logger.Panic().Logf("test panic")
	})
}

func TestStderrHandlerWithJsonFormatter(t *testing.T) {
	defer logging.Flush()

	logging.EnableCaller(logging.INFO)
	defer logging.DisableCaller(logging.INFO)
	logging.EnableCaller(logging.ERROR)
	defer logging.DisableCaller(logging.ERROR)

	logging.EnableStack(logging.ERROR)
	defer logging.DisableStack(logging.ERROR)

	h := NewStderrHandler(JsonFormatter)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	logging.Debug().Log()
	logging.Debug().Logf("stderr")

	logging.Info().Log()
	logging.Info().Logf("stderr")
	t.Run("parallel info", func(t *testing.T) {
		var w sync.WaitGroup
		for i := 0; i < 10; i++ {
			w.Add(1)
			go func(j int) {
				defer w.Done()

				logging.Info().Logf("test parallel")
			}(i)
		}
		w.Wait()
	})

	logging.Warn().Log()
	logging.Warn().Logf("error")

	logging.Error().Logf("error")

	logger := logging.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(logging.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(logging.ERROR)

	logger.Error().Logf("error")

	logger.EnableCaller(logging.PANIC)
	logger.EnableStack(logging.PANIC)
	assert.Panics(t, func() {
		logger.Panic().Logf("test panic")
	})
}

func TestStdoutHandlerWithLineFormatter(t *testing.T) {
	defer logging.Flush()

	logging.EnableCaller(logging.INFO)
	defer logging.DisableCaller(logging.INFO)
	logging.EnableCaller(logging.ERROR)
	defer logging.DisableCaller(logging.ERROR)

	logging.EnableStack(logging.ERROR)
	defer logging.DisableStack(logging.ERROR)

	h := NewStdoutHandler(StdFormatter)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	logging.Debug().Log()
	logging.Debug().Logf("stdout")

	logging.Info().Log()
	logging.Info().Logf("stdout")

	logging.Warn().Log()
	logging.Warn().Logf("error")

	logging.Error().Logf("error")

	logger := logging.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(logging.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(logging.ERROR)

	logger.Error().Logf("error")
}

func TestStdoutHandlerWithJsonFormatter(t *testing.T) {
	defer logging.Flush()

	logging.EnableCaller(logging.INFO)
	defer logging.DisableCaller(logging.INFO)
	logging.EnableCaller(logging.ERROR)
	defer logging.DisableCaller(logging.ERROR)

	logging.EnableStack(logging.ERROR)
	defer logging.DisableStack(logging.ERROR)

	h := NewStdoutHandler(JsonFormatter)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	logging.Debug().Log()
	logging.Debug().Logf("stdout")

	logging.Info().Log()
	logging.Info().Logf("stdout")

	logging.Warn().Log()
	logging.Warn().Logf("error")

	logging.Error().Logf("error")

	logger := logging.GetLogger("handler_test")
	defer logger.Flush()
	defer logger.Close()
	defer logger.SetPropagate(false)
	defer logger.DisableCaller(logging.ERROR)

	logger.SetPropagate(true)
	logger.EnableCaller(logging.ERROR)

	logger.Error().Logf("error")
}

func BenchmarkStderrHandler(b *testing.B) {
	defer logging.Flush()

	h := NewStderrHandler(StdFormatter)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logging.Info().Logf("stderr test")
		}
	})
}

func BenchmarkBufStderrHandler(b *testing.B) {
	defer logging.Flush()

	h, err := NewBufStderrHandler(StdFormatter)
	assert.Nil(b, err)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logging.Info().Logf("buf stderr test")
		}
	})
}

func BenchmarkStdoutHandler(b *testing.B) {
	defer logging.Flush()

	h := NewStdoutHandler(StdFormatter)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logging.Info().Logf("stdout test")
		}
	})
}

func BenchmarkBufStdoutHandler(b *testing.B) {
	defer logging.Flush()

	h, err := NewBufStdoutHandler(StdFormatter)
	assert.Nil(b, err)

	logging.AddHandler(h)
	defer logging.RemoveHandler(h)

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logging.Info().Logf("buf stdout test")
		}
	})
}
