package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// LogStdout indicates that logs should go to stdout
	LogStdout = "-"
)

// SimpleLogger is a simple logger that logs using JSON Lines.
type SimpleLogger struct {
	writer          io.Writer
	mutex           *sync.Mutex
	enableTimestamp bool
}

// NewSimpleLogger returns a new instance of SimpleLogger
func NewSimpleLogger(w io.Writer) *SimpleLogger {
	return &SimpleLogger{
		writer:          w,
		mutex:           &sync.Mutex{},
		enableTimestamp: true,
	}
}

// EnableTimestamp will enable logging of the timestamp
func (s *SimpleLogger) EnableTimestamp() {
	s.enableTimestamp = true
}

// DisableTimestamp will disable logging of the timestamp
func (s *SimpleLogger) DisableTimestamp() {
	s.enableTimestamp = false
}

// Marshal takes the interface and message and returns a byte array
func (s *SimpleLogger) Marshal(msg string, fields ...map[string]interface{}) ([]byte, error) {
	obj := map[string]interface{}{
		"msg": msg,
	}
	if s.enableTimestamp {
		obj["ts"] = time.Now().UTC().Format(time.RFC3339)
	}
	for _, f := range fields {
		for k, v := range f {
			obj[k] = v
		}
	}
	b, errMarshal := json.Marshal(obj)
	if errMarshal != nil {
		return make([]byte, 0), fmt.Errorf("error marshaling log entry: %w", errMarshal)
	}
	return b, nil
}

// Log is the interface to the simple logger to send a message
func (s *SimpleLogger) Log(msg string, fields ...map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	b, errMarshal := s.Marshal(msg, fields...)
	if errMarshal != nil {
		return errMarshal
	}
	_, errFprintln := fmt.Fprintln(s.writer, string(b))
	return errFprintln
}

// Error is the interface to log an error
func (s *SimpleLogger) Error(msg string, origErr error, fields ...map[string]interface{}) error {
	obj := map[string]interface{}{
		"error": origErr.Error(),
	}
	for _, f := range fields {
		for k, v := range f {
			obj[k] = v
		}
	}
	errLog := s.Log(msg, obj)
	if errLog != nil {
		return errLog
	}
	return nil
}

// Println prints a log message using fmt.Println without metadata
func (s *SimpleLogger) Println(input ...interface{}) {
	fmt.Println(input...)
}

// InitLogger initializes a logger to use stdout or a filename
func InitLogger(path string) (*SimpleLogger, error) {

	// Disable global flags from other packages
	log.SetPrefix("")
	log.SetFlags(0)

	switch path {
	case LogStdout:
		return NewSimpleLogger(os.Stdout), nil
	default:
		f, errOpenFile := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if errOpenFile != nil {
			return nil, fmt.Errorf("error opening log file %q: %w", path, errOpenFile)
		}
		return NewSimpleLogger(f), nil
	}
}

// DevNullLogger returns a logger that does nothing, useful for tests
func DevNullLogger() *SimpleLogger {
	return NewSimpleLogger(io.Discard)
}
