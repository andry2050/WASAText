package api

import (
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// uncommentMessage gestisce la rotta DELETE /messages/:message_id/comments/:comment_id
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Legge gli ID necessari dall'URL
	messageID := ps.ByName("message_id")
	commentID := ps.ByName("comment_id")

	// Rimuove il commento
	err := rt.db.UncommentMessage(messageID, commentID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la rimozione della reazione")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
