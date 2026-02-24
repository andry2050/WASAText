package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrUsernameInUse viene restituito quando si tenta di impostare un nome già preso
var ErrUsernameInUse = errors.New("username already in use")

// SetMyUserName aggiorna il nome utente se non è già preso da qualcun altro
func (db *appdbimpl) SetMyUserName(userID string, newName string) error {
	var existingUserID string

	err := db.c.QueryRow(`SELECT userid FROM users WHERE username = ?`, newName).Scan(&existingUserID)
	
	if err == nil {
		if existingUserID != userID {
			return ErrUsernameInUse // Nome già in uso
		}
		return nil 
		
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("errore controllo username: %w", err)
	}

	// Se il nome è disponibile, lo aggiorna
	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE userid = ?`, newName, userID)
	if err != nil {
		return fmt.Errorf("errore update username: %w", err)
	}

	return nil
}