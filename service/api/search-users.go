package api

import (
	"encoding/json"
	"net/http"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// searchUsers gestisce la rotta GET /users?username=...
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Controlla chi sta facendo la richiesta
	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Legge la richiesta dell'utente nella barra di ricerca
	searchQuery := r.URL.Query().Get("username")
	if searchQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ricerca nel database
	users, err := rt.db.SearchUsers(searchQuery)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la ricerca degli utenti")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Risponde inviando la lista JSON degli utenti trovati
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}