package database

import (
	"fmt"
)

// SetMyPhoto aggiorna la colonna photo_url dell'utente specificato.
func (db *appdbimpl) SetMyPhoto(userID string, photoPath string) error {
	// Fa un UPDATE per inserire il percorso dell'immagine
	res, err := db.c.Exec(`UPDATE users SET photo_url = ? WHERE userid = ?`, photoPath, userID)
	if err != nil {
		return fmt.Errorf("errore aggiornamento foto utente: %w", err)
	}

	// Verifica che l'utente esista davvero
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("impossibile verificare l'aggiornamento: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserDoesNotExist
	}

	return nil
}