package main

import "time"

import log "github.com/canhlinh/log4go"

const (
	filename = "flw.log"
)

func main() {
	// Get a new logger instance
	log.Close()
	// Create a default logger that is logging messages of FINE or higher

	/* Can also specify manually via the following: (these are the defaults) */
	flw := log.NewFileLogWriter(filename, false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(false)
	flw.SetRotateSize(0)
	flw.SetRotateLines(0)
	flw.SetRotateDaily(false)
	log.AddFilter("file", log.DEBUG, flw)
	// Log some experimental messages
	log.Finest("Everything is created now (notice that I will not be printing to the file)")
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
	log.Critical("Time to close out!")
	log.Close()
}
