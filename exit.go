package slog

import (
	"fmt"
	"os"
)

// global exit handler
var exitHandlers = make([]func(), 0)

func runExitHandlers() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Run exit handler error:", err)
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

// RegisterExitHandler register an exit-handler on global exitHandlers
func RegisterExitHandler(handler func()) {
	exitHandlers = append(exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func PrependExitHandler(handler func()) {
	exitHandlers = append([]func(){handler}, exitHandlers...)
}

// ResetExitHandlers reset all exitHandlers
func ResetExitHandlers(applyToStd bool) {
	exitHandlers = make([]func(), 0)

	if applyToStd {
		std.ResetExitHandlers()
	}
}

func (l *Logger) runExitHandlers() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Run exit handler error:", err)
		}
	}()

	for _, handler := range l.exitHandlers {
		handler()
	}
}

// RegisterExitHandler register an exit-handler on global exitHandlers
func (l *Logger) RegisterExitHandler(handler func()) {
	l.exitHandlers = append(l.exitHandlers, handler)
}

// PrependExitHandler prepend register an exit-handler on global exitHandlers
func (l *Logger) PrependExitHandler(handler func()) {
	l.exitHandlers = append([]func(){handler}, l.exitHandlers...)
}

// ResetExitHandlers reset logger exitHandlers
func (l *Logger) ResetExitHandlers() {
	l.exitHandlers = make([]func(), 0)
}

// ExitHandlers get all exitHandlers of the logger
func (l *Logger) ExitHandlers() []func() {
	return l.exitHandlers
}
