package badlogger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel int

const (
	None    LogLevel = -1
	Error   LogLevel = 0
	Warning LogLevel = 1
	Info    LogLevel = 2
	Develop LogLevel = 3
)

type BadLogger struct {
	Out             io.Writer
	TimeStampFormat string
	Level           LogLevel
}

func (bl *BadLogger) Write(b []byte) (n int, err error) { return (*bl).Out.Write(b) }

type stdWriter struct{}

func (s *stdWriter) Write(b []byte) (n int, err error) {
	fmt.Println(string(b))
	return len(b), nil
}

func DefaultLogger() BadLogger {
	return BadLogger{
		Out:             &stdWriter{},
		TimeStampFormat: time.Stamp,
		Level:           Develop,
	}
}

func FileLogger(logFilePath string, timeStampFormat string, logLevel LogLevel) (BadLogger, error) {

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return BadLogger{}, err
	}

	logger := BadLogger{
		Out:             logFile,
		TimeStampFormat: timeStampFormat,
		Level:           logLevel,
	}
	return logger, nil
}

func (bl *BadLogger) LogError(err error) {
	if bl.Level < Error {
		return
	}
	time := time.Now().Format(bl.TimeStampFormat)
	msg := fmt.Sprintf("%s :: %v", time, err)
	bl.Out.Write([]byte(msg))
}

func (bl *BadLogger) Log(msg string, level LogLevel) {
	if bl.Level < level {
		return
	}
	time := time.Now().Format(bl.TimeStampFormat)
	msg = fmt.Sprintf("%s :: %v", time, msg)
	bl.Out.Write([]byte(msg))

}
