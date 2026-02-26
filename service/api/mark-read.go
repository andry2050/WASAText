package api

import (
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) markMessagesAsRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	convID := ps.ByName("conversation_id")
	myID := extractBearer(r)

	err := rt.db.MarkMessagesAsRead(convID, myID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 OK, nessun contenuto
}
