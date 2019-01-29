// Package logger provides the global logger for the program
package logger

import "log"

// logger is the local variable for the logger
var logger *log.Logger

// Log provides the output function
func Log(v ...interface{}) {
	logger.Print(v...)
}

// Initialize sets the package up for logging
func Initialize(l *log.Logger) {
	logger = l
}
