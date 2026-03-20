package database

import (
	"fmt"
)

// LeaveGroup rimuove un utente dalla lista dei partecipanti di un gruppo
func (db *appdbimpl) LeaveGroup(groupID string, userID string) error {
	// Esegue la DELETE verificando che la conversazione sia di tipo 'group' evitando che l'utente abbandoni
	// per errore una chat singola.
	query := `
		DELETE FROM participants 
		WHERE userid = ? AND convid = ? 
		AND EXISTS (SELECT 1 FROM conversations WHERE convid = ? AND type = 'group')
	`

	res, err := db.c.Exec(query, userID, groupID, groupID)
	if err != nil {
		return fmt.Errorf("errore esecuzione query abbandono gruppo: %w", err)
	}

	// Controlla se è stato rimosso l'utente
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore controllo righe eliminate (abbandono gruppo): %w", err)
	}

	// Errore se l'utente non faceva parte del gruppo o il gruppo non esiste
	if rowsAffected == 0 {
		return ErrActionForbidden
	}

	return nil
}
