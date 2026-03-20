package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetConversation restituisce tutti i messaggi di una specifica conversazione
func (db *appdbimpl) GetConversation(convID string, userID string) (ConversationDetails, error) {
	// Controlla se l'utente fa parte della chat
	var isParticipant bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM participants WHERE convid = ? AND userid = ?)`, convID, userID).Scan(&isParticipant)
	if err != nil || !isParticipant {
		return ConversationDetails{}, fmt.Errorf("accesso negato o chat inesistente")
	}

	// Prende i messaggi con i dati del mittente
	query := `
		SELECT 
			m.msgid, m.senderid, u.username, u.photo_url, 
			m.content, m.is_photo, m.status, m.timestamp
		FROM messages m
		INNER JOIN users u ON m.senderid = u.userid
		WHERE m.convid = ?
		ORDER BY m.timestamp DESC
	`
	rows, err := db.c.Query(query, convID)
	if err != nil {
		return ConversationDetails{}, fmt.Errorf("errore lettura messaggi: %w", err)
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message
		var msgTime time.Time
		var photoURL sql.NullString

		err = rows.Scan(
			&msg.MessageID, &msg.SenderID, &msg.Sender.Username, &photoURL,
			&msg.Content, &msg.IsPhoto, &msg.Status, &msgTime,
		)
		if err != nil {
			return ConversationDetails{}, fmt.Errorf("errore scan messaggio: %w", err)
		}

		msg.Sender.UserID = msg.SenderID
		if photoURL.Valid {
			msg.Sender.PhotoURL = photoURL.String
		}
		msg.Timestamp = msgTime

		msg.Reactions = make([]Reaction, 0)

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return ConversationDetails{}, fmt.Errorf("errore iterazione messaggi: %w", err)
	}

	if messages == nil {
		messages = make([]Message, 0)
	}

	details := ConversationDetails{
		ConversationID: convID,
		Messages:       messages,
	}

	return details, nil
}
