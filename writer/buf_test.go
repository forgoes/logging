package writer

import "testing"

type stubWriter struct {
	writes  int
	flushed int
	closed  int
}

func (s *stubWriter) Write(p []byte) (int, error) {
	s.writes += len(p)
	return len(p), nil
}

func (s *stubWriter) Flush() error {
	s.flushed++
	return nil
}

func (s *stubWriter) Close() error {
	s.closed++
	return nil
}

func TestBufWriterCloseFlushesBeforeClose(t *testing.T) {
	w := &stubWriter{}
	bw, err := NewBufWriter(16, w)
	if err != nil {
		t.Fatalf("new buf writer: %v", err)
	}

	if _, err := bw.Write([]byte("hello")); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := bw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if w.flushed == 0 {
		t.Fatalf("expected flush to be called before close")
	}
	if w.closed == 0 {
		t.Fatalf("expected close to be called")
	}
}
