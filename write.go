package slog

//
// ---------------------------------------------------------------------------
// Do write log message
// ---------------------------------------------------------------------------
//

func (r *Record) logBytes(level Level) {
	// Will reduce memory allocation once
	// r.Message = strutil.Byte2str(message)

	// var buf *bytes.Buffer
	// buf = bufferPool.Get().(*bytes.Buffer)
	// defer bufferPool.Put(buf)
	// r.Buffer = buf

	// TODO release on here ??
	// defer r.logger.releaseRecord(r)

	handlers, ok := r.logger.matchHandlers(level)
	if !ok {
		return
	}

	// init record
	r.Level = level
	r.Init(r.logger.LowerLevelName)

	r.logger.mu.Lock()
	defer r.logger.mu.Unlock()

	// log caller. will alloc 3 times
	if r.logger.ReportCaller {
		caller, ok := getCaller(r.logger.CallerSkip)
		if ok {
			r.Caller = &caller
		}
	}

	// do write log message
	r.logger.write(level, r, handlers)

	// r.Buffer = nil
}

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

	if r.logger != nil {
		r.CallerFlag = r.logger.CallerFlag
	}

	r.microSecond = r.Time.Nanosecond() / 1000
}

//
// ---------------------------------------------------------------------------
// Do write log message
// ---------------------------------------------------------------------------
//

func (l *Logger) matchHandlers(level Level) ([]Handler, bool) {
	// alloc: 1 times for match handlers
	var matched []Handler
	for _, handler := range l.handlers {
		if handler.IsHandling(level) {
			matched = append(matched, handler)
		}
	}

	return matched, len(matched) > 0
}

func (l *Logger) write(level Level, r *Record, matched []Handler) {
	// // alloc: 1 times for match handlers
	// var matched []Handler
	// for _, handler := range l.handlers {
	// 	if handler.IsHandling(level) {
	// 		matched = append(matched, handler)
	// 	}
	// }
	//
	// // log level is don't match
	// if len(matched) == 0 {
	// 	return
	// }
	//
	// // init record
	// r.Init(l.LowerLevelName)
	// l.mu.Lock()
	// defer l.mu.Unlock()
	//
	// // log caller. will alloc 3 times
	// if l.ReportCaller {
	// 	caller, ok := getCaller(l.CallerSkip)
	// 	if ok {
	// 		r.Caller = &caller
	// 	}
	// }

	// processing log record
	for i := range l.processors {
		l.processors[i].Process(r)
	}

	// handling log record
	for _, handler := range matched {
		if err := handler.Handle(r); err != nil {
			printlnStderr("slog: failed to handle log, error: ", err)
		}
	}

	// flush logs on level <= error level.
	if level <= ErrorLevel {
		l.flushAll() // has locked on Record.logBytes()
	}

	if level <= PanicLevel {
		l.PanicFunc(r)
	} else if level <= FatalLevel {
		l.Exit(1)
	}
}
