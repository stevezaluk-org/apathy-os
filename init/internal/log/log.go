package log

import (
	"log"
	"os"
)

// New Initializes a new log.Logger and returns a pointer to it
func New(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix+": ", log.Lshortfile|log.LstdFlags)
}
