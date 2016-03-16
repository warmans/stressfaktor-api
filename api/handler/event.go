package handler

import (
	"net/http"
	"log"
	"github.com/warmans/stressfaktor-api/api/common"
	"strings"
	"strconv"
	"time"
	"github.com/warmans/stressfaktor-api/data/store"
	"golang.org/x/net/context"
)

func NewEventHandler(eventStore *store.Store) common.CtxHandler {
	return &EventHandler{eventStore: eventStore}
}

type EventHandler struct {
	eventStore *store.Store
}

func (h *EventHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context) {

	defer r.Body.Close()
	events, err := h.eventStore.FindEvents(h.filterFromRequest(r))
	if err != nil {
		log.Print(err.Error())
		common.SendResponse(rw, &common.Response{Status: http.StatusInternalServerError, Payload: nil})
		return
	}

	common.SendResponse(rw, &common.Response{Status: http.StatusOK, Payload: events})
}

func (h *EventHandler) filterFromRequest(r *http.Request) *store.EventFilter {
	r.ParseForm()

	filter := &store.EventFilter{
		EventIDs: make([]int, 0),
		VenueIDs: make([]int64, 0),
		Types: make([]string, 0),
		ShowDeleted: false,
	}

	if event := r.Form.Get("event"); event != "" {
		for _, idStr := range strings.Split(event, ",") {
			if idInt, err := strconv.Atoi(idStr); err == nil {
				filter.EventIDs = append(filter.EventIDs, idInt)
			}
		}
	}
	if from := r.Form.Get("from"); from != "" {
		if dateFrom, err := time.Parse("2006-01-02", from); err == nil {
			filter.DateFrom = dateFrom
		}
	}
	if to := r.Form.Get("to"); to != "" {
		if dateTo, err := time.Parse("2006-01-02", to); err == nil {
			filter.DateTo = dateTo
		}
	}
	if venue := r.Form.Get("venue"); venue != "" {
		for _, idStr := range strings.Split(venue, ",") {
			if idInt, err := strconv.Atoi(idStr); err == nil {
				filter.VenueIDs = append(filter.VenueIDs, int64(idInt))
			}
		}
	}
	if tpe := r.Form.Get("type"); tpe != "" {
		for _, typeStr := range strings.Split(tpe, ",") {
			filter.Types = append(filter.Types, typeStr)
		}
	}

	if deleted := r.Form.Get("deleted"); (deleted == "1" || deleted == "true") {
		filter.ShowDeleted = true
	}

	return filter
}
