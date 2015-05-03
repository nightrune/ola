package ola

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"
)

func TestMain(m *testing.M) {
	defer func() {
		if e := recover(); e != nil {
			logger.Debug(fmt.Sprintf("%s: %s", e, debug.Stack())) // line 20
		}
	}()
	logger.SetLoggingLevel(OLA_LOG_DEBUG)
	logger.Info("Running test for Ola Go Client Library... ")
	os.Exit(m.Run())
}
