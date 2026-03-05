package writer

import (
	"os"
)

type StdoutWriter struct {
}

func (s *StdoutWriter) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}

func (s *StdoutWriter) Flush() error {
	return os.Stdout.Sync()
}

func (s *StdoutWriter) Close() error {
	// Do not close process-wide stdout.
	return nil
}

func (s *StdoutWriter) WriteString(b string) (n int, err error) {
	return os.Stdout.WriteString(b)
}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

type StderrWriter struct {
}

func (s *StderrWriter) Write(b []byte) (n int, err error) {
	return os.Stderr.Write(b)
}

func (s *StderrWriter) Flush() error {
	return os.Stderr.Sync()
}

func (s *StderrWriter) Close() error {
	// Do not close process-wide stderr.
	return nil
}

func (s *StderrWriter) WriteString(b string) (n int, err error) {
	return os.Stderr.WriteString(b)
}

func NewStderrWriter() *StderrWriter {
	return &StderrWriter{}
}
