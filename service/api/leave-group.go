package api

import (
	"net/http"
	"strconv"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	groupid, err := strconv.Atoi(ps.ByName("groupid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	err = rt.db.leaveGroup(userid, groupid)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't leave the group")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Custom-Header", "Group leaved successfully")

}
