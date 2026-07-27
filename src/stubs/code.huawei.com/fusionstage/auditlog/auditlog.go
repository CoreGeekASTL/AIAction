package auditlog

import "io"

const DefaultUserID = 0
const DefaultGroupID = 0
const PERMISSION0440 = 0440

type Logger interface {
	RegisterSink(sink Sink)
	Print(content []byte) (int, error)
}

type Sink interface {
	Write(p []byte) (int, error)
}

type mockLogger struct {
	sinks []Sink
}

func NewLoggerBase(name string) Logger {
	return &mockLogger{sinks: make([]Sink, 0)}
}

func (l *mockLogger) RegisterSink(sink Sink) {
	l.sinks = append(l.sinks, sink)
}

func (l *mockLogger) Print(content []byte) (int, error) {
	for _, sink := range l.sinks {
		sink.Write(content)
	}
	return len(content), nil
}

func NewWriterSink(w io.Writer) Sink {
	return &writerSink{w: w}
}

type writerSink struct {
	w io.Writer
}

func (s *writerSink) Write(p []byte) (int, error) {
	return s.w.Write(p)
}