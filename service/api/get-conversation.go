package api

// TODO modificare tutti i return in http.error

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Prendo il userid e lo converto in intero
	userid, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	conversationid, err := strconv.Atoi(ps.ByName("conversationid"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Verifica se l'header Authorization è presente nella richiesta
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
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

	conversation := rt.db.GetConversation(userid, conversationid)
	if err != nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding the conversation")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
