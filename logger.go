package credui

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Mode controls how chatty the built-in logger is.
type Mode int

const (
	ModeProd Mode = iota
	ModeDev
)

// Logger is the logging interface used by this library.
//
// Libraries should never log unless explicitly configured by the user.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// nopLogger is the default logger.
// It produces no output and has near-zero overhead.
type nopLogger struct{}

func (nopLogger) Debugf(string, ...any) {}
func (nopLogger) Infof(string, ...any)  {}
func (nopLogger) Warnf(string, ...any)  {}
func (nopLogger) Errorf(string, ...any) {}

// stdLogger is an opt-in stdout logger intended for debugging and development.
type stdLogger struct {
	mu   sync.Mutex
	mode Mode
	l    *log.Logger
}

func newStdLogger(mode Mode) *stdLogger {
	// Timestamp is added manually to keep full control over the format.
	return &stdLogger{
		mode: mode,
		l:    log.New(os.Stdout, "", 0),
	}
}

func (s *stdLogger) Debugf(format string, args ...any) {
	if s.mode != ModeDev {
		return
	}
	s.write("DEBUG", format, args...)
}

func (s *stdLogger) Infof(format string, args ...any) {
	s.write("INFO", format, args...)
}

func (s *stdLogger) Warnf(format string, args ...any) {
	s.write("WARN", format, args...)
}

func (s *stdLogger) Errorf(format string, args ...any) {
	s.write("ERROR", format, args...)
}

func (s *stdLogger) write(level, format string, args ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, line := callerLocation(3)
	ts := time.Now().Format(time.RFC3339)

	msg := fmt.Sprintf(format, args...)
	msg = strings.ReplaceAll(msg, "\n", "\\n")

	s.l.Printf("%s [%s] %s:%d %s", ts, level, file, line, msg)
}

func callerLocation(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "?", 0
	}
	if idx := strings.LastIndex(file, string(os.PathSeparator)); idx >= 0 {
		file = file[idx+1:]
	}
	return file, line
}
