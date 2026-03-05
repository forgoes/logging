package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/forgoes/logging"
	"github.com/forgoes/logging/diode"
	"github.com/forgoes/logging/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateFileWriter struct {
	lj *lumberjack.Logger
}

func (l *RotateFileWriter) Write(p []byte) (int, error) {
	return l.lj.Write(p)
}

func (l *RotateFileWriter) Flush() error {
	return nil
}

func (l *RotateFileWriter) Close() error {
	return l.lj.Close()
}

func TestRingBufferHandlerWithLumberjackRotateFile(t *testing.T) {
	logDir := t.TempDir()
	logPath := filepath.Join(logDir, "app.log")

	rotate := &lumberjack.Logger{
		Filename: logPath,
		MaxSize:  1, // MB
		MaxAge:   7, // days
		Compress: false,
	}

	bw, err := writer.NewBufWriter(64*1024, &RotateFileWriter{lj: rotate})
	if err != nil {
		t.Fatalf("new buf writer: %v", err)
	}
	defer bw.Close()

	rh := NewRingBufferHandler(bw, JsonFormatter, 65536, diode.AlertFunc(func(int) {}), 0)

	logger := logging.GetLogger("handler.ring.lumberjack.test")
	logger.SetPropagate(false)
	logger.SetLevel(logging.DEBUG)
	logger.ResetHandler()
	logger.AddHandler(rh)
	defer logger.RemoveHandler(rh)
	defer rh.Close()

	// Keep enough volume to exceed MaxSize reliably, including debug mode.
	payload := strings.Repeat("X", 64*1024)
	for i := 0; i < 64; i++ {
		logger.Info().ALogf("line=%d payload=%s", i, payload)
	}

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		_ = rh.Flush()

		entries, readErr := os.ReadDir(logDir)
		if readErr != nil {
			t.Fatalf("read dir: %v", readErr)
		}

		rotated := 0
		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			name := e.Name()
			if strings.HasPrefix(name, "app-") && strings.HasSuffix(name, ".log") {
				rotated++
			}
		}

		if rotated > 0 {
			return
		}

		time.Sleep(50 * time.Millisecond)
	}

	entries, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatalf("read dir final: %v", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}

	t.Fatalf("rotation did not happen, files=%v", names)
}
