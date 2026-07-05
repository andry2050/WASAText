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

	// Analizza il form multipart
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Legge il testo dal FormData di Vue (chiave "content")
	content := r.FormValue("content")
	var photoURL string

	// Controlla se c'è un'immagine allegata (chiave "image")
	file, _, errFile := r.FormFile("image")
	if errFile == nil {
		defer file.Close()

		photoUUID, _ := uuid.NewV4()
		photoFileName := photoUUID.String() + ".jpg"

		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		photoPath := filepath.Join("uploads", photoFileName)
		dst, err := os.Create(photoPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		photoURL = "/" + photoPath
	}

	// Deve esserci ALMENO il testo O la foto
	if content == "" && photoURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Invia al database passando entrambi
	msg, err := rt.db.SendMessage(convID, senderID, content, photoURL)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(msg)
}
