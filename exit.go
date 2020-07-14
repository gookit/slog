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

// PrependExitHandler register an exit-handler on global exitHandlers
func RegisterExitHandler(handler func())  {
	exitHandlers = append(exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func PrependExitHandler(handler func())  {
	exitHandlers = append([]func(){handler}, exitHandlers...)
}

func (logger *Logger) runExitHandlers()  {
	defer func() {
		if err := recover(); err != nil {
			_,_ = fmt.Fprintln(os.Stderr, "Run exit handler error:", err)
		}
	}()

	for _, handler := range exitHandlers {
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
