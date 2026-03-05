package handler

import (
	"sync"

	"github.com/forgoes/logging"
	"github.com/forgoes/logging/writer"
)

type StderrHandler struct {
	sync.Mutex
	format Formatter
	w      *writer.StderrWriter
}

func NewStderrHandler(format Formatter) *StderrHandler {
	return &StderrHandler{
		format: format,
		w:      writer.NewStderrWriter(),
	}
}

func (s *StderrHandler) Handle(e *logging.Event) (bool, error) {
	b, err := s.format(e)
	if err != nil {
		return true, err
	}

	s.Lock()
	defer s.Unlock()

	if _, err := s.w.Write(b); err != nil {
		return true, err
	}

	if _, err := s.w.WriteString("\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (s *StderrHandler) Flush() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Flush()
}

func (s *StderrHandler) Close() error {
	s.Lock()
	defer s.Unlock()

	// Do not close process-wide stderr.
	return s.w.Flush()
}

type BufStderrHandler struct {
	sync.Mutex
	format Formatter
	w      *writer.BufWriter
}

func NewBufStderrHandler(format Formatter) (*BufStderrHandler, error) {
	w, err := writer.NewBufWriter(0, writer.NewStderrWriter())
	if err != nil {
		return nil, err
	}

	return &BufStderrHandler{
		format: format,
		w:      w,
	}, nil
}

func (s *BufStderrHandler) Handle(e *logging.Event) (bool, error) {
	b, err := s.format(e)
	if err != nil {
		return true, err
	}

	s.Lock()
	defer s.Unlock()

	if _, err := s.w.Write(b); err != nil {
		return true, err
	}

	if _, err := s.w.WriteString("\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (s *BufStderrHandler) Flush() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Flush()
}

func (s *BufStderrHandler) Close() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Close()
}

type StdoutHandler struct {
	sync.Mutex
	format Formatter
	w      *writer.StdoutWriter
}

func NewStdoutHandler(format Formatter) *StdoutHandler {
	return &StdoutHandler{
		format: format,
		w:      writer.NewStdoutWriter(),
	}
}

func (s *StdoutHandler) Handle(e *logging.Event) (bool, error) {
	b, err := s.format(e)
	if err != nil {
		return true, err
	}

	s.Lock()
	defer s.Unlock()

	if _, err := s.w.Write(b); err != nil {
		return true, err
	}

	if _, err := s.w.WriteString("\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (s *StdoutHandler) Flush() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Flush()
}

func (s *StdoutHandler) Close() error {
	s.Lock()
	defer s.Unlock()

	// Do not close process-wide stdout.
	return s.w.Flush()
}

type BufStdoutHandler struct {
	sync.Mutex
	format Formatter
	w      *writer.BufWriter
}

func NewBufStdoutHandler(format Formatter) (*BufStdoutHandler, error) {
	w, err := writer.NewBufWriter(0, writer.NewStdoutWriter())
	if err != nil {
		return nil, err
	}

	return &BufStdoutHandler{
		format: format,
		w:      w,
	}, nil
}

func (s *BufStdoutHandler) Handle(e *logging.Event) (bool, error) {
	b, err := s.format(e)
	if err != nil {
		return true, err
	}

	s.Lock()
	defer s.Unlock()

	if _, err := s.w.Write(b); err != nil {
		return true, err
	}

	if _, err := s.w.WriteString("\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (s *BufStdoutHandler) Flush() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Flush()
}

func (s *BufStdoutHandler) Close() error {
	s.Lock()
	defer s.Unlock()

	return s.w.Close()
}
