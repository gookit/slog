package slog

import "context"

// SubLogger is a sub-logger, It can be used to keep a certain amount of contextual information and log multiple times.
// 可以用于保持一定的上下文信息多次记录日志。例如在循环中使用，或者作为方法参数传入。
//
// Usage:
//
//	sl := slog.NewSub().KeepCtx(custom ctx).
//		KeepFields(slog.M{"ip": ...}).
//		KeepData(slog.M{"username": ...})
//	defer sl.Release()
//
//	sl.Info("some message")
type SubLogger struct {
	l *Logger // parent logger

	// Ctx keep context for all log records
	Ctx context.Context
	// Fields keep custom fields data for all log records
	Fields M
	// Data keep data for all log records
	Data M
	// Extra data. will keep for all log records
	Extra M
}

// NewSubWith returns a new SubLogger with parent logger.
func NewSubWith(l *Logger) *SubLogger { return &SubLogger{l: l} }

// KeepCtx keep context for all log records
func (sub *SubLogger) KeepCtx(ctx context.Context) *SubLogger {
	sub.Ctx = ctx
	return sub
}

// KeepFields keep custom fields data for all log records
func (sub *SubLogger) KeepFields(fields M) *SubLogger {
	sub.Fields = fields
	return sub
}

// KeepField keep custom field for all log records
func (sub *SubLogger) KeepField(field string, value any) *SubLogger {
	if sub.Fields == nil {
		sub.Fields = make(M)
	}

	sub.Fields[field] = value
	return sub
}

// KeepData keep data for all log records
func (sub *SubLogger) KeepData(data M) *SubLogger {
	sub.Data = data
	return sub
}

// KeepExtra keep extra data for all log records
func (sub *SubLogger) KeepExtra(extra M) *SubLogger {
	sub.Extra = extra
	return sub
}

// Release releases the SubLogger.
func (sub *SubLogger) Release() {
	sub.l = nil
	sub.Ctx = nil
	sub.Fields = nil
	sub.Data = nil
	sub.Extra = nil
}

func (sub *SubLogger) withKeepCtx() *Record {
	r := sub.l.WithContext(sub.Ctx)
	r.Data = sub.Data
	r.Extra = sub.Extra
	r.Fields = sub.Fields
	return r
}

//
// ---------------------------------------------------------------------------
// Add log message with level
// ---------------------------------------------------------------------------
//

// Print logs a message at PrintLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Print(args ...any) { sub.withKeepCtx().Print(args...) }

// Printf logs a message at PrintLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Printf(format string, args ...any) { sub.withKeepCtx().Printf(format, args...) }

// Trace logs a message at TraceLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Trace(args ...any) { sub.withKeepCtx().Trace(args...) }

// Tracef logs a formatted message at TraceLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Tracef(format string, args ...any) {
	sub.withKeepCtx().Tracef(format, args...)
}

// Debug logs a message at DebugLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Debug(args ...any) { sub.withKeepCtx().Debug(args...) }

// Debugf logs a formatted message at DebugLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Debugf(format string, args ...any) {
	sub.withKeepCtx().Debugf(format, args...)
}

// Info logs a message at InfoLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Info(args ...any) { sub.withKeepCtx().Info(args...) }

// Infof logs a formatted message at InfoLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Infof(format string, args ...any) {
	sub.withKeepCtx().Infof(format, args...)
}

// Notice logs a message at NoticeLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Notice(args ...any) { sub.withKeepCtx().Notice(args...) }

// Noticef logs a formatted message at NoticeLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Noticef(format string, args ...any) {
	sub.withKeepCtx().Noticef(format, args...)
}

// Warn logs a message at WarnLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Warn(args ...any) { sub.withKeepCtx().Warn(args...) }

// Warnf logs a formatted message at WarnLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Warnf(format string, args ...any) {
	sub.withKeepCtx().Warnf(format, args...)
}

// Error logs a message at ErrorLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Error(args ...any) { sub.withKeepCtx().Error(args...) }

// Errorf logs a formatted message at ErrorLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Errorf(format string, args ...any) {
	sub.withKeepCtx().Errorf(format, args...)
}

// Fatal logs a message at FatalLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Fatal(args ...any) { sub.withKeepCtx().Fatal(args...) }

// Fatalf logs a formatted message at FatalLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Fatalf(format string, args ...any) {
	sub.withKeepCtx().Fatalf(format, args...)
}

// Panic logs a message at PanicLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Panic(args ...any) { sub.withKeepCtx().Panic(args...) }

// Panicf logs a formatted message at PanicLevel. will with sub logger's context, fields and data
func (sub *SubLogger) Panicf(format string, args ...any) {
	sub.withKeepCtx().Panicf(format, args...)
}
