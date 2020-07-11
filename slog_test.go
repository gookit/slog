package slog

import "testing"

func TestInfof(t *testing.T) {
	AddProcessor()

	Infof("info %s", "message")
}
