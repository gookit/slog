package main

import (
	log "github.com/gookit/slog"
)

const simplestTemplate = "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}"

func init() {
	log.GetFormatter().(*log.TextFormatter).SetTemplate(simplestTemplate)
	log.SetLogLevel(log.ErrorLevel)
	log.Errorf("Test")
}

func main() {
}
