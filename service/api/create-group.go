package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// createGroup gestisce la rotta POST /groups
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	creatorID := extractBearer(r)
	if creatorID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Legge il nome del gruppo e i membri dal JSON
	var reqBody struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupName := strings.TrimSpace(reqBody.Name)
	if groupName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newGroup, err := rt.db.CreateGroup(groupName, reqBody.Members, creatorID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la creazione del gruppo")
		// Potrebbe esserci un ID utente non valido nella lista dei membri
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newGroup)
}