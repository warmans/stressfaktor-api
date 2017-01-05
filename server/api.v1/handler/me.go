package handler

import (
	"net/http"

	"github.com/warmans/fakt-api/server/api.v1/common"
	"github.com/warmans/fakt-api/server/api.v1/middleware"
)

func NewMeHandler() http.Handler {
	return &MeHandler{}
}

type MeHandler struct{}

func (h *MeHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	common.SendResponse(
		rw,
		&common.Response{
			Status:  http.StatusOK,
			Payload: middleware.GetUser(r),
		},
	)
}
