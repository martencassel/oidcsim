package middleware

import (
	"context"
	"net/http"

	"github.com/martencassel/oidcsim/internal/application/session"
)

type ctxKeySession struct{}

var sessionKey = ctxKeySession{}

func WithSessionManager(mgr session.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ensure a session exists (or create one)
			sid, err := mgr.Ensure(w, r)
			if err != nil {
				http.Error(w, "session error", http.StatusInternalServerError)
				return
			}
			// Store both the manager and the current session ID in context
			ctx := context.WithValue(r.Context(), sessionKey, mgr)
			ctx = context.WithValue(ctx, "sessionID", sid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SessionManagerFromContext(ctx context.Context) (session.SessionManager, bool) {
	mgr, ok := ctx.Value(sessionKey).(session.SessionManager)
	return mgr, ok
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	sid, ok := ctx.Value("sessionID").(string)
	return sid, ok
}
