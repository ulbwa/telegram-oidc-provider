package common

import "context"

// contextKey is a type for context keys to avoid collisions
type contextKey string

const (
	// TraceIDKey is the context key for trace ID
	TraceIDKey contextKey = "trace_id"
)

// GetTraceID extracts trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}
