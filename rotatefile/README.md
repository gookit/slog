# Rotate File

 `rotatefile` provides simple file rotation, compression and cleanup.

## Features

- Rotate file by size and time
  - Custom filename for rotate file by size
  - Custom time clock for rotate
  - Custom file perm for create log file
  - Custom rotate mode: create, rename
- Compress rotated file
- Cleanup old files

## Install

```bash
go get github.com/gookit/slog/rotatefile
```

## Usage

### Create a file writer

```go
logFile := "testdata/go_logger.log"
writer, err := rotatefile.NewConfig(logFile).Create()
if err != nil {
    panic(err)
}

// use writer
writer.Write([]byte("log message\n"))
```

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

### Available config options

```go
// Config struct for rotate dispatcher
type Config struct {
    // Filepath the log file path, will be rotating
    Filepath string `json:"filepath" yaml:"filepath"`
    
    // FilePerm for create log file. default DefaultFilePerm
    FilePerm os.FileMode `json:"file_perm" yaml:"file_perm"`
    
    // MaxSize file contents max size, unit is bytes.
    // If is equals zero, disable rotate file by size
    //
    // default see DefaultMaxSize
    MaxSize uint64 `json:"max_size" yaml:"max_size"`
    
    // RotateTime the file rotate interval time, unit is seconds.
    // If is equals zero, disable rotate file by time
    //
    // default see EveryHour
    RotateTime RotateTime `json:"rotate_time" yaml:"rotate_time"`
    
    // CloseLock use sync lock on write contents, rotating file.
    //
    // default: false
    CloseLock bool `json:"close_lock" yaml:"close_lock"`
    
    // BackupNum max number for keep old files.
    //
    // 0 is not limit, default is DefaultBackNum
    BackupNum uint `json:"backup_num" yaml:"backup_num"`
    
    // BackupTime max time for keep old files, unit is hours.
    //
    // 0 is not limit, default is DefaultBackTime
    BackupTime uint `json:"backup_time" yaml:"backup_time"`
    
    // Compress determines if the rotated log files should be compressed using gzip.
    // The default is not to perform compression.
    Compress bool `json:"compress" yaml:"compress"`
    
    // RenameFunc you can custom-build filename for rotate file by size.
    //
    // default see DefaultFilenameFn
    RenameFunc func(filePath string, rotateNum uint) string
    
    // TimeClock for rotate
    TimeClock Clocker
}
```

## Files clear

```go
	fc := rotatefile.NewFilesClear(func(c *rotatefile.CConfig) {
		c.AddPattern("/path/to/some*.log")
		c.BackupNum = 2
		c.BackupTime = 12 // 12 hours
	})
	
	// clear files on daemon
	go fc.DaemonClean()
	
	// NOTE: stop daemon before exit
	// fc.QuitDaemon()
```

### Configs

```go

```