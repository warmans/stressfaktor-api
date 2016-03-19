package handler

import (
	"net/http"
	"github.com/warmans/stressfaktor-api/api/common"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

func NewLogoutHandler(sess sessions.Store) common.CtxHandler {
	return &LogoutHandler{sessions: sess}
}

type LogoutHandler struct {
	sessions sessions.Store
}

func (h *LogoutHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context) {
	defer r.Body.Close()
	r.ParseForm()

	//create their session
	sess, err := h.sessions.Get(r, "login")
	if err == nil {
		delete(sess.Values, "login")
		sess.Options.MaxAge = -1;
		sess.Save(r, rw)
	}
	common.SendResponse(
		rw,
		&common.Response{
			Status: http.StatusOK,
			Payload: nil,
		},
	)
}