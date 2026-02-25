package database

import (
	"fmt"
)

// SetGroupPhoto aggiorna la foto di un gruppo, verificando che il richiedente sia un membro
func (db *appdbimpl) SetGroupPhoto(groupID string, photoPath string, requesterID string) error {
	query := `
		UPDATE conversations 
		SET photo_url = ? 
		WHERE convid = ? AND type = 'group' 
		AND EXISTS (SELECT 1 FROM participants WHERE convid = ? AND userid = ?)
	`
	
	res, err := db.c.Exec(query, photoPath, groupID, groupID, requesterID)
	if err != nil {
		return fmt.Errorf("errore esecuzione query aggiornamento foto gruppo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore controllo righe aggiornate: %w", err)
	}

	if rowsAffected == 0 {
		return ErrActionForbidden
	}

	return nil
}