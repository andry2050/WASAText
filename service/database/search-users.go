package database

import (
	"database/sql"
	"fmt"
)

// SearchUsers cerca nel database gli utenti il cui username contiene la stringa specificata
func (db *appdbimpl) SearchUsers(searchQuery string) ([]User, error) {
	likeQuery := "%" + searchQuery + "%"

	// Esegue la ricerca
	rows, err := db.c.Query(`SELECT userid, username, photo_url FROM users WHERE username LIKE ?`, likeQuery)
	if err != nil {
		return nil, fmt.Errorf("errore query ricerca utenti: %w", err)
	}
	defer rows.Close() 

	var users []User
	
	// Scorre i risultati trovati
	for rows.Next() {
		var u User
		// La foto profilo potrebbe non esserci, quindi uso sql.NullString per evitare che il programma vada in crash
		var photoURL sql.NullString 
		
		err = rows.Scan(&u.UserID, &u.Username, &photoURL)
		if err != nil {
			return nil, fmt.Errorf("errore lettura riga utente: %w", err)
		}
		
		// Se la foto c'è la assegna all'utente
		if photoURL.Valid {
			u.PhotoURL = photoURL.String
		}
		
		// Aggiunge l'utente alla lista
		users = append(users, u)
	}

	// Controlla che non ci siano stati errori durante lo scorrimento
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("errore iterazione righe database: %w", err)
	}

	if users == nil {
		users = make([]User, 0)
	}

	return users, nil
}