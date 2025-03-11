package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {


	userid, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupid, err := strconv.Atoi(ps.ByName("groupid"))
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
	if groupid != authToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var updatedGroup Group
	err = json.NewDecoder(r.Body).Decode(&updatedGroup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !updatedGroup.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedGroup.GroupID = groupid
	updatedGroup.UserID = userid

	err = rt.db.UpdateGroupName(updatedUser.ToDatabase())
	if errors.Is(err, database.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).WithField("groupid", groupid).Error("can't update the group name")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
