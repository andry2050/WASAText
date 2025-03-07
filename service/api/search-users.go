package api

import (
	"encoding/json"
	"regexp"
	"strings"

	"net/http"
	"strconv"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Verifica se l'header Authorization è presente nella richiesta
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token di autorizzazione mancante", http.StatusUnauthorized)
		return
	}

	// Estrai il token dall'header Authorization
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 1 {
		http.Error(w, "Formato del token non valido", http.StatusUnauthorized)
		return
	}

	authToken := tokenParts[0]

	authTokenCast, err := strconv.Atoi(authToken)
	if err != nil {

		http.Error(w, "Token di autorizzazione not uint64", http.StatusBadRequest)
		return
	}

	userid := authTokenCast

	// Get the search query from the request
	query_search := r.URL.Query().Get("search")
	validQuerySearch, err := regexp.MatchString("^.*?$", query_search)
	if err != nil {
		if !(validQuerySearch && 3 <= len(query_search) && len(query_search) <= 16) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	users, err := rt.db.SearchUsers(userid, query_search)
	if err != nil {
		ctx.Logger.Error("Error searching users ", err)
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	searchedUsers := make([]User, len(users))
	for i, u := range users {
		var user User
		user.FromDatabase(u)
		searchedUsers[i] = user
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchedUsers); err != nil {
		ctx.Logger.Error("Error encoding users ", err)
		http.Error(w, "Error encoding response ", http.StatusInternalServerError)
		return
	}

}
