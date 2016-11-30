package middleware

import (
	"net/http"

	"context"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/sessions"
	"github.com/warmans/fakt-api/server/data/service/user"
)

func AddCtx(nextHandler http.Handler, sess sessions.Store, users *user.UserStore, logger log.Logger) http.Handler {
	return &CtxMiddleware{next: nextHandler, sessions: sess, users: users, logger: logger}
}

type CtxMiddleware struct {
	next     http.Handler
	sessions sessions.Store
	users    *user.UserStore
	logger   log.Logger
}

func (m *CtxMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	//setup logger with http context info and add to context
	logger := log.NewContext(m.logger).With("method", r.Method, "url", r.URL.String())
	ctx = context.WithValue(ctx, "logger", m.logger)

	sess, err := m.sessions.Get(r, "login")
	if err != nil {
		logger.Log("msg", "Failed to get session: "+err.Error())
		m.next.ServeHTTP(rw, r.WithContext(ctx))
		return
	}

	userId, found := sess.Values["userId"]
	if found == false || userId == nil || userId.(int64) < 1 {
		m.next.ServeHTTP(rw, r.WithContext(ctx))
		return
	}

	usr, err := m.users.GetUser(userId.(int64))
	if err == nil && usr != nil {
		ctx = context.WithValue(ctx, "user", usr)
	} else {
		ctx = context.WithValue(ctx, "user", nil)
	}

	m.next.ServeHTTP(rw, r.WithContext(ctx))
}
