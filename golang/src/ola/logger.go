package ola

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "Log: ", log.Ldate|log.Ltime|log.Lshortfile)
}
