package api

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"

	"github.com/Viron35/wasa/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var photo Photo

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

	err = r.ParseMultipartForm(10 << 20) // Max size 10MB
	if err != nil {
		http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	// Leggi il file dell'immagine dall'oggetto Request
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Errore nel leggere il file dell'immagine", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Legge i dati dell'immagine in un array di byte
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

	// Copia i dati dal file ricevuto al file di destinazione
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Errore nel salvataggio del file", http.StatusInternalServerError)
		return
	}

	imageURL := "http://yourserver.com/uploads/" + filepath.Base(savePath)

	// Popola il struct `photo` con l'URL dell'immagine invece di base64
	photo.UserID = userid
	photo.Image = imageURL



	err = rt.db.UpdatePhoto(photo.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't upload the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Custom-Header", "Photo loaded successfully")

}
