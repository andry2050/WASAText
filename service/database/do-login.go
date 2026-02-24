package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// DoLogin verifica se un utente esiste: se esiste restituisce l'ID, altrimenti lo crea
func (db *appdbimpl) DoLogin(username string) (string, error) {
	var userID string
	
	// Cerca se l'utente esiste già
	err := db.c.QueryRow(`SELECT userid FROM users WHERE username = ?`, username).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			newUUID, errUUID := uuid.NewV4()
			if errUUID != nil {
				return "", fmt.Errorf("errore generazione uuid: %w", errUUID)
			}
			userID = newUUID.String()

			_, errInsert := db.c.Exec(`INSERT INTO users (userid, username) VALUES (?, ?)`, userID, username)
			if errInsert != nil {
				return "", fmt.Errorf("errore inserimento nuovo utente: %w", errInsert)
			}
			
			return userID, nil
		}

		return "", err
	}

	// Se l'utente già esisteva restituisce il suo ID
	return userID, nil
}