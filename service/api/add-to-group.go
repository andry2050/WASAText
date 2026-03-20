package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// addToGroup gestisce la rotta POST /groups/:group_id/members
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requesterID := extractBearer(r)
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("group_id")

	// Legge l'ID dell'utente da aggiungere dal JSON
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	targetUserID := reqBody.UserID
	if targetUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Passa al database il gruppo, chi viene aggiunto, e chi lo sta aggiungendo
	err := rt.db.AddToGroup(groupID, targetUserID, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'aggiunta al gruppo")
		// Se l'errore è per permessi negati (es. non fa parte del gruppo)
		if errors.Is(err, database.ErrActionForbidden) {
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden
			return
		}
		// Altrimenti errore generico (es. utente non esiste)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
