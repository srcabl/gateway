package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

// InjectSession handles injecting the ResponseWriter and Request structs
// into context so that resolver methods can use these to set and read cookies
func InjectSession(session *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpContext := HTTP{
				W: &w,
				R: r,
			}

			ctx := context.WithValue(r.Context(), HTTPKey, httpContext)
			ctx = context.WithValue(ctx, SessionKey, session)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// CtxKey is the key used to extract from the context
type CtxKey string

const (
	// HTTPKey is the key used to extract the http struct
	HTTPKey CtxKey = "http"
	//SessionKey is the key used to extract the session
	SessionKey CtxKey = "session"
)

// HTTP is the struct used to inject the response writer and request http structs
type HTTP struct {
	W *http.ResponseWriter
	R *http.Request
}

// GetSession returns a cached session of the given name
func GetSession(ctx context.Context, name string) *sessions.Session {
	store := ctx.Value(SessionKey).(*sessions.CookieStore)
	httpContext := ctx.Value(HTTPKey).(HTTP)

	// Ignore err because a session is always returned even if one doesn't exist
	session, _ := store.Get(httpContext.R, name)

	return session
}

// SaveSession saves the session by writing it to the response
func SaveSession(ctx context.Context, session *sessions.Session) error {
	httpContext := ctx.Value(HTTPKey).(HTTP)
	if err := session.Save(httpContext.R, *httpContext.W); err != nil {
		return errors.Wrap(err, "failed to save session")
	}
	return nil
}
