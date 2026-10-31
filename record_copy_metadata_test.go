package slog_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gookit/goutil/x/assert"
	"github.com/gookit/slog"
)

type copyContextKey struct{}

func TestRecordCopyPreservesMetadata(t *testing.T) {
	ctx := context.WithValue(context.Background(), copyContextKey{}, "request-42")
	stamp := time.Date(2024, time.March, 2, 12, 30, 0, 123000000, time.UTC)
	cases := []struct {
		name   string
		derive func(*slog.Record) *slog.Record
	}{
		{"Copy", func(r *slog.Record) *slog.Record { return r.Copy() }},
		{"WithField", func(r *slog.Record) *slog.Record { return r.WithField("name", "value") }},
		{"WithFields", func(r *slog.Record) *slog.Record { return r.WithFields(slog.M{"name": "value"}) }},
		{"WithData", func(r *slog.Record) *slog.Record { return r.WithData(slog.M{"name": "value"}) }},
		{"WithError", func(r *slog.Record) *slog.Record { return r.WithError(errors.New("example")) }},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			l := slog.New()
			r := l.Record().SetContext(ctx)
			r.Time = stamp
			derived := tt.derive(r)
			assert.True(t, derived.Ctx == ctx)
			assert.True(t, derived.Time.Equal(stamp))
			assert.True(t, r.Ctx == ctx)
			assert.True(t, r.Time.Equal(stamp))
		})
	}
}

func TestRecordCopyWithContextAndTime(t *testing.T) {
	ctx := context.WithValue(context.Background(), copyContextKey{}, "original")
	replacement := context.WithValue(context.Background(), copyContextKey{}, "replacement")
	stamp := time.Date(2024, time.March, 2, 12, 30, 0, 0, time.UTC)
	r := slog.New().Record().SetContext(ctx)
	r.Time = stamp
	t.Run("WithContext preserves time", func(t *testing.T) {
		n := r.WithContext(replacement)
		assert.True(t, n.Ctx == replacement)
		assert.True(t, n.Time.Equal(stamp))
	})
	t.Run("WithTime preserves context", func(t *testing.T) {
		n := r.WithTime(stamp.Add(time.Hour))
		assert.True(t, n.Ctx == ctx)
		assert.True(t, n.Time.Equal(stamp.Add(time.Hour)))
	})
}
