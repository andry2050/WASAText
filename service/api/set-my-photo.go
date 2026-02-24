package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// setMyPhoto gestisce la rotta PUT /users/me/photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := extractBearer(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Limita la dimensione del file (es. massimo 10 MB) per non intasare il server
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.Logger.WithError(err).Error("File troppo grande o formattato male")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Estrae il file dalla richiesta
	file, _, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.WithError(err).Error("Nessun file trovato nella richiesta")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crea un nome unico per questa foto usando un UUID
	photoUUID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photoFileName := photoUUID.String() + ".jpg" 
	
	// Crea la cartella "uploads" se non esiste già
	os.MkdirAll("uploads", os.ModePerm)
	photoPath := filepath.Join("uploads", photoFileName)

	// Crea il file vuoto sul disco del server
	dst, err := os.Create(photoPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Impossibile salvare il file sul disco")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copia i dati arrivati da internet dentro il file vuoto appena creato
	if _, err := io.Copy(dst, file); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la copia dei dati del file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Il database aggiorna il profilo dell'utente con il nuovo percorso della foto
	dbPhotoPath := "/" + photoPath 
	err = rt.db.SetMyPhoto(userID, dbPhotoPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore database in setMyPhoto")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}