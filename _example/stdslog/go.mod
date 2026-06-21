// Separate module: imports log/slog (Go 1.21+). Uses the standalone
// github.com/gookit/rotatefile as an io.Writer for std slog.
module stdslog-example

go 1.21

require github.com/gookit/rotatefile v0.1.0

require (
	github.com/gookit/goutil v0.7.6 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)
