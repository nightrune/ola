/**
 * @brief This class is a simple wrapper for adding a few use things to go
 * logging, and help keep similarity between the other Ola Client libraries and
 * this one
 *
 * It provides a global logger, that is initialzed on init.
 *
 * Right now its a simple wrapper over logging but intended to get added to
 * later, and allow a consistent api
 */
package ola

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
)

var logger *OlaLogger

func init() {
	_log := log.New(os.Stdout, "", 0)
	logger = NewLogger(_log)
}

const (
	OLA_LOG_NONE  = iota
	OLA_LOG_FATAL = iota
	OLA_LOG_WARN  = iota
	OLA_LOG_INFO  = iota
	OLA_LOG_DEBUG = iota
)

type OlaLogger struct {
	logger    *log.Logger
	log_level uint
}

func NewLogger(log_interface *log.Logger) *OlaLogger {
	l := new(OlaLogger)
	l.logger = log_interface
	l.log_level = OLA_LOG_NONE
	return l
}

func (m *OlaLogger) SetLoggingLevel(level uint) {
	if level > OLA_LOG_DEBUG {
		return
	}
	m.log_level = level
}

func (m *OlaLogger) log(log_type uint, msg string) {
	if log_type > m.log_level {
		return
	}

	if m.logger != nil {
		_, file, line, ok := runtime.Caller(2)
		if ok == false {
			m.logger.Print(msg)
			return
		}
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("%s:%d: ", file, line))
		buffer.WriteString(msg)
		m.logger.Print(buffer.String())
	}
}

func (m *OlaLogger) Debug(msg string) {
	m.log(OLA_LOG_DEBUG, msg)
}

func (m *OlaLogger) Warn(msg string) {
	m.log(OLA_LOG_WARN, msg)
}

func (m *OlaLogger) Fatal(msg string) {
	m.log(OLA_LOG_FATAL, msg)
}

func (m *OlaLogger) Info(msg string) {
	m.log(OLA_LOG_INFO, msg)
}
