package log

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLogger(t *testing.T) {
	// Create the logger
	stdoutLogger, errInitLogger := InitLogger(LogStdout)
	assert.NoError(t, errInitLogger)
	stdoutLogger.EnableTimestamp()

	// Log to a buff that can be introspected
	buf := bytes.NewBuffer([]byte{})
	logger := NewSimpleLogger(buf)
	logger.DisableTimestamp()

	// Log a message with no fields
	errLog := logger.Log("Log message")
	assert.NoError(t, errLog)
	assert.Equal(t, "{\"msg\":\"Log message\"}\n", buf.String())
	buf.Reset()

	// Log a message with fields
	errLog = logger.Log("Log message", map[string]interface{}{
		"key1": "value1",
	})
	assert.NoError(t, errLog)
	assert.Equal(t, "{\"key1\":\"value1\",\"msg\":\"Log message\"}\n", buf.String())
	buf.Reset()

	errLog = logger.Log("Log message", map[string]interface{}{
		"key1": "value1",
	}, map[string]interface{}{
		"key2": "value2",
	})
	assert.NoError(t, errLog)
	assert.Equal(t, "{\"key1\":\"value1\",\"key2\":\"value2\",\"msg\":\"Log message\"}\n", buf.String())
	buf.Reset()

	// Log an Error
	errError := logger.Error("this is an error", errors.New("original error"))
	assert.NoError(t, errError)
	assert.Equal(t, "{\"error\":\"original error\",\"msg\":\"this is an error\"}\n", buf.String())
	buf.Reset()

	// Log an Error with fields
	errError = logger.Error("this is an error", errors.New("original error"), map[string]interface{}{
		"key1": "value1",
	}, map[string]interface{}{
		"key2": "value2",
	})
	assert.NoError(t, errError)
	assert.Equal(t, "{\"error\":\"original error\",\"key1\":\"value1\",\"key2\":\"value2\",\"msg\":\"this is an error\"}\n", buf.String())
	buf.Reset()

	// Log as Println
	logger.Println("test")

	// Use devnull logger
	devNull := DevNullLogger()
	errLog = devNull.Log("goes nowhere")
	assert.NoError(t, errLog)

	// Init logger with filename
	filename := "delete.log"
	defer os.Remove(filename)
	fileLog, errInitLogger := InitLogger(filename)
	fileLog.DisableTimestamp()
	assert.NoError(t, errInitLogger)
	errLog = fileLog.Log("to file")
	assert.NoError(t, errLog)

	f, errOpen := os.Open(filename)
	assert.NoError(t, errOpen)
	defer f.Close()

	out, errReadAll := ioutil.ReadAll(f)
	assert.NoError(t, errReadAll)
	assert.Equal(t, "{\"msg\":\"to file\"}\n", string(out))
}
