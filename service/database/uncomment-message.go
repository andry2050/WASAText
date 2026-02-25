package database

import (
	"fmt"
)

// UncommentMessage rimuove una reazione precedentemente inserita
func (db *appdbimpl) UncommentMessage(messageID string, reactionID string, userID string) error {
	// Esegue la DELETE verificando contemporaneamente l'ID del messaggio, della reazione e dell'utente
	// Questo impedisce agli utenti di cancellare le reazioni degli altri
	res, err := db.c.Exec(`DELETE FROM reactions WHERE reactionid = ? AND msgid = ? AND userid = ?`, reactionID, messageID, userID)
	if err != nil {
		return fmt.Errorf("errore esecuzione query di eliminazione reazione: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore controllo righe eliminate (reazione): %w", err)
	}

	// Errore se la reazione non esiste o l'utente sta cercando di cancellare quella di un altro
	if rowsAffected == 0 {
		return fmt.Errorf("nessuna reazione eliminata (permesso negato o inesistente)")
	}

	return nil
}