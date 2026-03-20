package database

import (
	"errors"
	"fmt"
)

// Definiamo un errore specifico per quando l'utente prova a cancellare un messaggio non suo
var ErrMessageNotFoundOrForbidden = errors.New("messaggio non trovato o non autorizzato")

// DeleteMessage elimina un messaggio dal database, solo se chi lo richiede ne è il mittente
func (db *appdbimpl) DeleteMessage(messageID string, userID string) error {
	// Esegue la DELETE mettendo come condizione sia l'ID del messaggio, sia l'ID del mittente
	// In questo modo è letteralmente impossibile che un utente cancelli il messaggio di un altro
	query := `DELETE FROM messages WHERE msgid = ? AND senderid = ?`

	res, err := db.c.Exec(query, messageID, userID)
	if err != nil {
		return fmt.Errorf("errore esecuzione query di eliminazione: %w", err)
	}

	// Controlla quante righe sono state effettivamente cancellate
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("errore controllo righe eliminate: %w", err)
	}

	// Se non è stata cancellata nessuna riga, significa che il messaggio non esiste
	// oppure l'utente stava provando a cancellare un messaggio altrui
	if rowsAffected == 0 {
		return ErrMessageNotFoundOrForbidden
	}

	return nil
}
