package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLogin gestisce la rotta POST /session
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Struttura per leggere il body della richiesta
	var reqBody struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Impossibile decodificare il body JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(reqBody.Name)

	// Validazione: tra 3 e 16 caratteri altrimenti errore
	if len(username) < 3 || len(username) > 16 {
		ctx.Logger.Warning("Username non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	identifier, err := rt.db.DoLogin(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore nel database durante il login")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		Identifier string `json:"identifier"`
	}{
		Identifier: identifier,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
