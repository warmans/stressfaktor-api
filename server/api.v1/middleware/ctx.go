package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/warmans/ctxhandler"
	"github.com/warmans/fakt-api/server/api.v1/common"
	"github.com/warmans/fakt-api/server/data/store"
	"golang.org/x/net/context"
)

func AddCtx(nextHandler ctxhandler.CtxHandler, sess sessions.Store, users *store.UserStore, restrict bool) http.Handler {
	return &CtxMiddleware{next: nextHandler, sessions: sess, users: users, restrict: restrict}
}

type CtxMiddleware struct {
	next     ctxhandler.CtxHandler
	sessions sessions.Store
	users    *store.UserStore
	restrict bool
}

func (m *CtxMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	sess, err := m.sessions.Get(r, "login")
	if err != nil {
		log.Printf("Failed to get session: %s", err.Error())
		m.next.ServeHTTP(rw, r, ctx)
		return
	}

	userId, found := sess.Values["userId"]
	if found == false || userId == nil || userId.(int64) < 1 {
		if m.restrict {
			common.SendError(rw, common.HTTPError{"Access Denied", http.StatusForbidden, nil}, false)
			return
		}
		m.next.ServeHTTP(rw, r, ctx)
		return
	}

	user, err := m.users.GetUser(userId.(int64))
	if err == nil && user != nil {
		ctx = context.WithValue(ctx, "user", user)
	} else {
		if m.restrict {
			common.SendError(rw, common.HTTPError{"Access Denied", http.StatusForbidden, nil}, false)
			return
		}
	}

	m.next.ServeHTTP(rw, r, ctx)
}
