package database

import (
	"errors"
	"fmt"
)

// ErrActionForbidden viene restituito quando un utente prova a fare un'azione su un gruppo a cui non appartiene
var ErrActionForbidden = errors.New("azione non consentita: non fai parte del gruppo")

// AddToGroup aggiunge un utente a un gruppo, verificando che il richiedente sia già un membro.
func (db *appdbimpl) AddToGroup(groupID string, targetUserID string, requesterID string) error {
	// Verifica che la conversazione sia di tipo 'group' e che il requester ne faccia parte
	var isAuthorized bool
	authQuery := `
		SELECT EXISTS(
			SELECT 1 FROM conversations c
			INNER JOIN participants p ON c.convid = p.convid
			WHERE c.convid = ? AND c.type = 'group' AND p.userid = ?
		)
	`
	err := db.c.QueryRow(authQuery, groupID, requesterID).Scan(&isAuthorized)
	if err != nil || !isAuthorized {
		return ErrActionForbidden
	}

	// Inserisce il nuovo utente nel gruppo usando INSERT OR IGNORE così se l'utente è già nel gruppo
	// il database non va in errore e non duplica i dati
	insertQuery := `INSERT OR IGNORE INTO participants (convid, userid) VALUES (?, ?)`
	_, err = db.c.Exec(insertQuery, groupID, targetUserID)
	if err != nil {
		return fmt.Errorf("errore inserimento nuovo partecipante: %w", err)
	}

	return nil
}