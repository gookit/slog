package slog

import (
	"fmt"
	"os"
)

// global exit handler
var exitHandlers = make([]func(), 0)

func runExitHandlers()  {
	defer func() {
		if err := recover(); err != nil {
			_,_ = fmt.Fprintln(os.Stderr, "Run exit handler error:", err)
		}
	}()

	for _, handler := range exitHandlers {
		handler()
	}
}

// ExitHandlers get all global exitHandlers
func ExitHandlers() []func() {
	return exitHandlers
}

// PrependExitHandler register an exit-handler on global exitHandlers
func RegisterExitHandler(handler func())  {
	exitHandlers = append(exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func PrependExitHandler(handler func())  {
	exitHandlers = append([]func(){handler}, exitHandlers...)
}

// ResetExitHandlers reset all exitHandlers
func ResetExitHandlers(applyToStd bool)  {
	exitHandlers = make([]func(), 0)

	if applyToStd {
		std.ResetExitHandlers()
	}
}

func (logger *Logger) runExitHandlers()  {
	defer func() {
		if err := recover(); err != nil {
			_,_ = fmt.Fprintln(os.Stderr, "Run exit handler error:", err)
		}
	}()

	for _, handler := range logger.exitHandlers {
		handler()
	}
}

// PrependExitHandler register an exit-handler on global exitHandlers
func (logger *Logger) RegisterExitHandler(handler func())  {
	logger.exitHandlers = append(logger.exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func (logger *Logger) PrependExitHandler(handler func())  {
	logger.exitHandlers = append([]func(){handler}, logger.exitHandlers...)
}

// ResetExitHandlers reset logger exitHandlers
func (logger *Logger) ResetExitHandlers()  {
	logger.exitHandlers = make([]func(), 0)
}

// ExitHandlers get all exitHandlers of the logger
func (logger *Logger) ExitHandlers() []func() {
	return logger.exitHandlers
}