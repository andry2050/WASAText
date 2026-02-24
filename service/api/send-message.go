package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// sendMessage gestisce la rotta POST /conversations/:conversation_id/messages
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	senderID := extractBearer(r)
	if senderID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversation_id")

	// 3. Analizzia il form multipart (limite 10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore nel parsing del form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var content string
	var isPhoto bool

	// Controlle se è una foto o un testo
	file, _, errFile := r.FormFile("file")
	if errFile == nil {

		defer file.Close()
		isPhoto = true

		// Salva l'immagine
		photoUUID, _ := uuid.NewV4()
		photoFileName := photoUUID.String() + ".jpg"
		os.MkdirAll("uploads", os.ModePerm)
		photoPath := filepath.Join("uploads", photoFileName)

		dst, err := os.Create(photoPath)
		if err != nil {
			ctx.Logger.WithError(err).Error("Impossibile salvare la foto del messaggio")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		content = "/" + photoPath 
	} else {	
		isPhoto = false
		content = r.FormValue("text")
		
		// Se non c'è né foto né testo, la richiesta non è valida
		if content == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Invia al database
	msg, err := rt.db.SendMessage(convID, senderID, content, isPhoto)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore salvataggio messaggio nel DB")
		// Potrebbe essere che l'utente stia provando a scrivere in una chat a cui non appartiene
		w.WriteHeader(http.StatusForbidden) 
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(msg)
}