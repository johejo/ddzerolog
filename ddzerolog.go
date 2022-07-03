// Package ddzerolog provides a log/span correlation func for the github.com/rs/zerolog.
package ddzerolog

import (
	"context"

	"github.com/rs/zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// UpdateContext updates the zerolog internal logger's context to attach trace and span details found in the given context.
func UpdateContext(ctx context.Context) func(c zerolog.Context) zerolog.Context {
	return func(c zerolog.Context) zerolog.Context {
		span, ok := tracer.SpanFromContext(ctx)
		if !ok {
			return c
		}
		return c.Uint64("dd.trace_id", span.Context().TraceID()).Uint64("dd.span_id", span.Context().SpanID())
	}
}
