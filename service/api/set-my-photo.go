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

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var photo Photo

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
	if groupid != authToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // Max size 10MB
	if err != nil {
		http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Errore nel leggere il file dell'immagine", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	photoData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Errore nella lettura dei dati dell'immagine", http.StatusInternalServerError)
		return
	}

	savePath := "./uploads/" + strconv.Itoa(userid) + "_" + time.Now().Format("20060102150405") + ".jpg"

	// Crea il file di destinazione
	outFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Errore nella creazione del file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()


	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !updatedUser.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedUser.UserID = userid
	updatedUser.Photo = imageURL

	err = rt.db.SetMyPhoto(updatedUser.ToDatabase())
	if errors.Is(err, database.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).WithField("userid", userid).Error("can't update your photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
