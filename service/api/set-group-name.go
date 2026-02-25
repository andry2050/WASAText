package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// setGroupName gestisce la rotta PUT /groups/:group_id/name
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requesterID := extractBearer(r)
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("group_id")

	var reqBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newName := strings.TrimSpace(reqBody.Name)
	if newName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := rt.db.SetGroupName(groupID, newName, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la modifica del nome del gruppo")
		if errors.Is(err, database.ErrActionForbidden) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}