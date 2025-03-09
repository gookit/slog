package internal

import "path/filepath"

// AddSuffix2path add suffix to file path.
//
// eg: "/path/to/error.log" => "/path/to/error.{suffix}.log"
func AddSuffix2path(filePath, suffix string) string {
	ext := filepath.Ext(filePath)
	return filePath[:len(filePath)-len(ext)] + "." + suffix + ext
}

// BuildGlobPattern builds a glob pattern for the given logfile. NOTE: use for testing only.
func BuildGlobPattern(logfile string) string {
	return logfile[:len(logfile)-4] + "*"
}
