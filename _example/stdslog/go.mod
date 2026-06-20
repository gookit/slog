// Separate module: this example imports log/slog (needs Go 1.21+), while the
// main gookit/slog module stays on go 1.20. Kept under _example (ignored by the
// main module's ./... builds) so it doesn't affect the library's Go support.
module stdslog-example

go 1.21

require github.com/gookit/slog v0.0.0

require (
	github.com/gookit/goutil v0.7.6 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/gookit/slog => ../..
