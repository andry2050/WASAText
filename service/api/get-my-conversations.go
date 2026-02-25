package api

import (
	"encoding/json"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyConversations gestisce la rotta GET /conversations
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Conversazioni dell'utente
	conversations, err := rt.db.GetMyConversations(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante il recupero delle conversazioni")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conversations)
}