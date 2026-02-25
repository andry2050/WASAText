package api

import (
	"encoding/json"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getConversation gestisce la rotta GET /conversations/:conversation_id
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversation_id")

	convDetails, err := rt.db.GetConversation(conversationID, userID)
	if err != nil {
		// Se c'è un errore, potrebbe essere che la chat non esiste o l'utente non partecipa
		ctx.Logger.WithError(err).Error("Errore durante il recupero della conversazione")
		w.WriteHeader(http.StatusNotFound) 
		return
	}

	// Invia i dati al frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(convDetails)
}