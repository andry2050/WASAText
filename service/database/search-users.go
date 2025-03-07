package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) SearchUsers(userid int, search string) ([]User, error) {

	var users []User

	rows, err := db.c.Query(
		`SELECT userid, username FROM users WHERE username LIKE '%' || ? || '%' ORDER BY username `, search)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var u User
		if err := rows.Scan(&u.UserID, &u.Username); err != nil {
			return nil, err
		}

/*		isBanned, err := db.CheckBan(u.UserID, userid)
		if err != nil {
			return nil, err
		}
		if !isBanned {
			users = append(users, u)
		}
*/
	}

	defer func() { err = rows.Close() }() // Chiudi le righe del risultato alla fine della funzione

	return users, err
}
