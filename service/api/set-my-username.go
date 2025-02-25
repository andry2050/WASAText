package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/andry2050/wasa/service/api/reqcontext"
	"github.com/andry2050/wasa/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Prendo il userid e lo converto in intero
	userid, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verifica se l'header Authorization è presente nella richiesta
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	authToken, err := strconv.Atoi(authHeader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Valida il token
	if userid != authToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !updatedUser.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedUser.UserID = userid

	err = rt.db.UpdateUsername(updatedUser.ToDatabase())
	if errors.Is(err, database.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).WithField("userid", userid).Error("can't update the username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
