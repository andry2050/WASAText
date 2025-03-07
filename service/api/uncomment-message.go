package api

import (
	"net/http"
	"strconv"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var comment Comment

	commentid, err := strconv.Atoi(ps.ByName("commentid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Prendo il userid e lo converto in intero
	userid, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postid, err := strconv.Atoi(ps.ByName("messageid"))
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

	comment.CommentID = commentid
	comment.UserID = authToken
	comment.PostID = postid

	err = rt.db.UncommentMessage(comment.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't uncomment the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Custom-Header", "Photo uncomment successfully")
	

}
