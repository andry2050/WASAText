package api

import (
	"encoding/json"
	"net/http"
	"github.com/andry2050/wasa/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !user.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exist, err := rt.db.UserExists(user.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't check if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exist {
		dbuser, err := rt.db.GetUser(user.Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't get the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.FromDatabase(dbuser)

	} else {
		dbuser, err := rt.db.CreateUser(user.ToDatabase())
		if err != nil {
			ctx.Logger.WithError(err).Error("can't create the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.FromDatabase(dbuser)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
