package recovery

import (
	"log/slog"
	"runtime/debug"
)

// Handle recovers from a panic and logs it with a stack trace.
// Use as: defer recovery.Handle()
func Handle() {
	if r := recover(); r != nil {
		slog.Warn("Recovered panic", "panic", r, "trace", string(debug.Stack()))
	}
}
