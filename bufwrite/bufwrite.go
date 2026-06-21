// Package bufwrite is kept for backward compatibility only.
//
// Deprecated: bufwrite has moved to its own module. Use
// github.com/gookit/rotatefile/bufwrite instead. This package only re-exports
// the new location and will be removed in a future major version.
package bufwrite

import "github.com/gookit/rotatefile/bufwrite"

// BufIOWriter is an alias.
//
// Deprecated: use github.com/gookit/rotatefile/bufwrite.BufIOWriter
type BufIOWriter = bufwrite.BufIOWriter

// LineWriter is an alias.
//
// Deprecated: use github.com/gookit/rotatefile/bufwrite.LineWriter
type LineWriter = bufwrite.LineWriter

// Deprecated: use github.com/gookit/rotatefile/bufwrite.NewBufIOWriter
var NewBufIOWriter = bufwrite.NewBufIOWriter

// Deprecated: use github.com/gookit/rotatefile/bufwrite.NewBufIOWriterSize
var NewBufIOWriterSize = bufwrite.NewBufIOWriterSize

// Deprecated: use github.com/gookit/rotatefile/bufwrite.NewLineWriter
var NewLineWriter = bufwrite.NewLineWriter

// Deprecated: use github.com/gookit/rotatefile/bufwrite.NewLineWriterSize
var NewLineWriterSize = bufwrite.NewLineWriterSize
