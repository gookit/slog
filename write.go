package slog

//
// ---------------------------------------------------------------------------
// Do write log message
// ---------------------------------------------------------------------------
//

// func (r *Record) logWrite(level Level) {
// Will reduce memory allocation once
// r.Message = strutil.Byte2str(message)

// var buf *bytes.Buffer
// buf = bufferPool.Get().(*bytes.Buffer)
// defer bufferPool.Put(buf)
// r.Buffer = buf

// TODO release on here ??
// defer r.logger.releaseRecord(r)
// r.logger.writeRecord(level, r)
// r.Buffer = nil
// }

// Init something for record.
func (r *Record) Init(lowerLevelName bool) {
	// use lower level name
	if lowerLevelName {
		r.levelName = r.Level.LowerName()
	} else {
		r.levelName = r.Level.Name()
	}

	// init log time
	if r.Time.IsZero() {
		r.Time = r.logger.TimeClock.Now()
	}

	// r.microSecond = r.Time.Nanosecond() / 1000
}

// Init something for record.
func (r *Record) beforeHandle(l *Logger) {
	// log caller. will alloc 3 times
	if l.ReportCaller {
		caller, ok := getCaller(r.CallerSkip)
		if ok {
			r.Caller = &caller
		}
	}

	// processing log record
	for i := range l.processors {
		l.processors[i].Process(r)
	}
}

// do write log record
func (l *Logger) writeRecord(level Level, r *Record) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// do write log message
	var inited bool
	for _, handler := range l.handlers {
		if handler.IsHandling(level) {
			if !inited {
				// init, call processors
				r.Init(l.LowerLevelName)
				r.beforeHandle(l)
				inited = true
			}

			if err := handler.Handle(r); err != nil {
				l.err = err
				printlnStderr("slog: failed to handle log, error:", err)
			}
		}
	}

	// ---- after write log ----

	// flush logs on level <= error level.
	if level <= ErrorLevel {
		l.flushAll() // has been in lock
	}

	if level <= PanicLevel {
		l.PanicFunc(r)
	} else if level <= FatalLevel {
		l.Exit(1)
	}
}
