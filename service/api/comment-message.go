package api

import (
	"encoding/json"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// commentMessage gestisce la rotta POST /messages/:message_id/comments
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	messageID := ps.ByName("message_id")

	// Legge l'emoticon dal JSON inviato dal frontend
	var reqBody struct {
		Emoji string `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil || reqBody.Emoji == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Salva la reazione nel database
	reaction, err := rt.db.CommentMessage(messageID, userID, reqBody.Emoji)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'aggiunta della reazione")
		// L'utente potrebbe star provando a reagire a un messaggio di una chat a cui non appartiene
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(reaction)
}
