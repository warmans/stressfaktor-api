package middleware

import (
	"net/http"

	"github.com/warmans/fakt-api/pkg/server/api.v1/common"
	"github.com/go-kit/kit/log"
)

func AddSetup(nextHandler http.Handler) http.Handler {
	return &SetupMiddleware{next: nextHandler}
}

type SetupMiddleware struct {
	next    http.Handler
	headers map[string]string
	logger log.Logger
}

func (m *SetupMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {


	//always parse the form
	if err := r.ParseForm(); err != nil {
		common.SendError(rw, common.HTTPError{"Invalid request data", http.StatusBadRequest, nil}, nil)
		return
	}
	//always close the body
	defer func() {
		if err := r.Body.Close(); err != nil {
			//can't really do much once the headers have already been sent. Meh!
			m.logger.Log("msg", "Failed to close HTTP response body:", err)
		}
	}()
	m.next.ServeHTTP(rw, r)
}