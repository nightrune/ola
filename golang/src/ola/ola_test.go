package ola

import (
	"log"
	"os"
	"runtime/debug"
	"testing"
)

func TestMain(m *testing.M) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack()) // line 20
		}
	}()
	os.Exit(m.Run())
}
