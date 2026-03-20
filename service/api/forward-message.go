package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// forwardMessage gestisce la rotta POST /conversations/:conversation_id/forward
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	senderID := extractBearer(r)
	if senderID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	targetConvID := ps.ByName("conversation_id")

	// Legge l'ID del messaggio da inoltrare dal JSON
	var reqBody struct {
		MessageID string `json:"message_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqBody.MessageID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Chiama il database per eseguire l'inoltro
	newMsg, err := rt.db.ForwardMessage(reqBody.MessageID, targetConvID, senderID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'inoltro del messaggio")
		// Se l'errore è dovuto a permessi o messaggio non trovato
		if errors.Is(err, database.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newMsg)
}
