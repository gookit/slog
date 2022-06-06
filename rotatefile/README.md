# Rotate File

## Usage

### Use on another logger

```go
package main

import (
  "log"

  "github.com/gookit/slog/rotatefile"
)

func main() {
	logFile := "testdata/go_logger.log"
	writer, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err) 
	}

	log.SetOutput(writer)
	log.Println("log message")
}
```
