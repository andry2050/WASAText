package api

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andry2050/WASAText/service/api/reqcontext"
	"github.com/andry2050/WASAText/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// setGroupPhoto gestisce la rotta PUT /groups/:group_id/photo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requesterID := extractBearer(r)
	if requesterID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("group_id")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoUUID, _ := uuid.NewV4()
	photoFileName := photoUUID.String() + ".jpg"
	os.MkdirAll("uploads", os.ModePerm)
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

	dbPhotoPath := "/" + photoPath 
	err = rt.db.SetGroupPhoto(groupID, dbPhotoPath, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante l'aggiornamento della foto del gruppo")
		if errors.Is(err, database.ErrActionForbidden) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}