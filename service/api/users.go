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

// setMyUserName gestisce la rotta PUT /users/me/username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrae l'ID dell'utente dall'header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := strings.TrimPrefix(authHeader, "Bearer ")

	// Legge il nuovo nome dal JSON inviato
	var reqBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newName := strings.TrimSpace(reqBody.Name)

	// Verifica del nuovo nome
	if len(newName) < 3 || len(newName) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Chiama il database
	err := rt.db.SetMyUserName(userID, newName)
	if err != nil {
		// Controlla se l'username è già stato usato, dando l'errore 409
		if errors.Is(err, database.ErrUsernameInUse) {
			w.WriteHeader(http.StatusConflict) // 409 Conflict
			return
		}
		// Altrimenti è un errore generico del server
		ctx.Logger.WithError(err).Error("Errore del database in setMyUserName")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
