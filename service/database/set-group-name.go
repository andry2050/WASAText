package database

import (
	"fmt"
)

// SetGroupName aggiorna il nome di un gruppo, verificando che il richiedente sia un membro
func (db *appdbimpl) SetGroupName(groupID string, newName string, requesterID string) error {
	// Esegue l'UPDATE solo se il gruppo esiste e l'utente ne fa parte
	query := `
		UPDATE conversations 
		SET name = ? 
		WHERE convid = ? AND type = 'group' 
		AND EXISTS (SELECT 1 FROM participants WHERE convid = ? AND userid = ?)
	`
	
	res, err := db.c.Exec(query, newName, groupID, groupID, requesterID)
	if err != nil {
		return fmt.Errorf("errore esecuzione query aggiornamento nome gruppo: %w", err)
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