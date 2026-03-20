package api

import (
	"errors"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// leaveGroup gestisce la rotta DELETE /groups/:group_id/members/:user_id
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requesterID := extractBearer(r)
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("group_id")
	targetUserID := ps.ByName("user_id")

	if requesterID != targetUserID {
		ctx.Logger.Warning("Un utente ha provato a rimuovere un'altra persona dal gruppo")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Rimozione dell'utente dal gruppo
	err := rt.db.LeaveGroup(groupID, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'abbandono del gruppo")
		// Errore se il gruppo non esiste o l'utente non partecipa
		if errors.Is(err, database.ErrActionForbidden) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
