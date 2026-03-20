package database

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
)

// CommentMessage aggiunge un'emoticon come reazione a un messaggio.
func (db *appdbimpl) CommentMessage(messageID string, userID string, emoji string) (Reaction, error) {
	// Controlla se l'utente fa parte della chat in cui è stato inviato il messaggio
	var isParticipant bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM messages m
			INNER JOIN participants p ON m.convid = p.convid
			WHERE m.msgid = ? AND p.userid = ?
		)
	`
	err := db.c.QueryRow(checkQuery, messageID, userID).Scan(&isParticipant)
	if err != nil || !isParticipant {
		return Reaction{}, fmt.Errorf("messaggio non trovato o accesso negato")
	}

	reactionUUID, _ := uuid.NewV4()
	reactionID := reactionUUID.String()

	// Inserisce la reazione nel database
	_, err = db.c.Exec(`INSERT INTO reactions (reactionid, msgid, userid, emoji) VALUES (?, ?, ?, ?)`, reactionID, messageID, userID, emoji)
	if err != nil {
		return Reaction{}, fmt.Errorf("errore inserimento reazione: %w", err)
	}

	// Recupera le informazioni dell'utente
	var u User
	var photoURL sql.NullString
	_ = db.c.QueryRow(`SELECT username, photo_url FROM users WHERE userid = ?`, userID).Scan(&u.Username, &photoURL)
	u.UserID = userID
	if photoURL.Valid {
		u.PhotoURL = photoURL.String
	}

	return Reaction{
		ReactionID: reactionID,
		MessageID:  messageID,
		UserID:     userID,
		User:       u,
		Emoji:      emoji,
	}, nil
}
