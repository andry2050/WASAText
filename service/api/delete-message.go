package api

import (
	"errors"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// deleteMessage gestisce la rotta DELETE /messages/:message_id
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	messageID := ps.ByName("message_id")

	// Eliminazione dal database
	err := rt.db.DeleteMessage(messageID, userID)
	if err != nil {
		// Se l'errore è dovuto al fatto che il messaggio non è suo o non esiste
		if errors.Is(err, database.ErrMessageNotFoundOrForbidden) {
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden
			return
		}

		ctx.Logger.WithError(err).Error("Errore durante l'eliminazione del messaggio")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
