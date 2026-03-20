package database

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
)

// CreateGroup crea una nuova conversazione di tipo "group" e aggiunge i partecipanti
func (db *appdbimpl) CreateGroup(groupName string, memberIDs []string, creatorID string) (Group, error) {

	groupUUID, _ := uuid.NewV4()
	groupID := groupUUID.String()

	// Unisce il creatore alla lista dei membri, evitando duplicati
	uniqueMembers := make(map[string]bool)
	uniqueMembers[creatorID] = true
	for _, id := range memberIDs {
		uniqueMembers[id] = true
	}

	// Se uno degli inserimenti fallisce, annulla tutto per non lasciare dati a metà
	tx, err := db.c.Begin()
	if err != nil {
		return Group{}, fmt.Errorf("impossibile avviare transazione: %w", err)
	}

	// Crea la conversazione principale
	_, err = tx.Exec(`INSERT INTO conversations (convid, type, name, photo_url) VALUES (?, 'group', ?, '')`, groupID, groupName)
	if err != nil {
		tx.Rollback()
		return Group{}, fmt.Errorf("errore inserimento gruppo: %w", err)
	}

	// Inserisce tutti i membri nella tabella participanti
	for userID := range uniqueMembers {
		_, err = tx.Exec(`INSERT INTO participants (convid, userid) VALUES (?, ?)`, groupID, userID)
		if err != nil {
			tx.Rollback()
			// Nel caso in cui l'ID di un utente inviato non esiste genera errore
			return Group{}, fmt.Errorf("errore inserimento partecipante %s: %w", userID, err)
		}
	}

	// Salva la transazione
	if err = tx.Commit(); err != nil {
		return Group{}, fmt.Errorf("impossibile confermare transazione: %w", err)
	}

	// Costruisce l'oggetto di risposta recuperando i dati completi degli utenti appena inseriti
	var members []User
	for userID := range uniqueMembers {
		var u User
		var photoURL sql.NullString
		// Prendendo nome e foto di ogni partecipante
		err := db.c.QueryRow(`SELECT username, photo_url FROM users WHERE userid = ?`, userID).Scan(&u.Username, &photoURL)
		if err == nil {
			u.UserID = userID
			if photoURL.Valid {
				u.PhotoURL = photoURL.String
			}
			members = append(members, u)
		}
	}

	newGroup := Group{
		GroupID:  groupID,
		Name:     groupName,
		PhotoURL: "",
		Members:  members,
	}

	return newGroup, nil
}
