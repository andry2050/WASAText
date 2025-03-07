package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var comment Comment

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Errore nella decodifica JSON", http.StatusBadRequest)
		return
	}

	// Estrae il testo dal JSON
	text, ok := jsonData["text"].(string)
	if !ok {
		http.Error(w, "Errore: il valore associato a 'text' non è una stringa.", http.StatusBadRequest)
		return
	}

	// Prendo il userid e lo converto in intero
	userid, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postid, err := strconv.Atoi(ps.ByName("postid"))
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

	// Conversione da stringa a intero
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Valida il token
//	check, _ := rt.db.CheckBan(userid, authToken)


	comment.UserID = authToken
	comment.MessageID = messageid
	comment.Text = text

	err = rt.db.CommentPhoto(comment.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't comment the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Custom-Header", "Comment posted successfully")

}
