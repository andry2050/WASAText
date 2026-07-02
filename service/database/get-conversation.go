package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) GetConversation(targetOrConvID string, userID string) (ConversationDetails, error) {
	var realConvID string
	var isConv bool

	// Controlla se l'ID passato è di una conversazione esistente
	db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM conversations WHERE convid = ?)`, targetOrConvID).Scan(&isConv)

	if isConv {
		realConvID = targetOrConvID
	} else {
		// Se non è una conversazione, calcola l'ID univoco della chat diretta
		if userID < targetOrConvID {
			realConvID = userID + "_" + targetOrConvID
		} else {
			realConvID = targetOrConvID + "_" + userID
		}
	}

	// Controlla se l'utente fa parte della chat
	var isParticipant bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM participants WHERE convid = ? AND userid = ?)`, realConvID, userID).Scan(&isParticipant)

	if err != nil || !isParticipant {
		// La chat non esiste ancora perché non ci sono messaggi.
		// Restituisce una chat vuota senza dare errore
		return ConversationDetails{
			ConversationID: realConvID,
			Messages:       make([]Message, 0),
		}, nil
	}

	// Prende i messaggi
	query := `
		SELECT 
			m.msgid, m.senderid, u.username, u.photo_url, 
			m.content, m.is_photo, m.status, m.timestamp
		FROM messages m
		INNER JOIN users u ON m.senderid = u.userid
		WHERE m.convid = ?
		ORDER BY m.timestamp DESC
	`
	rows, err := db.c.Query(query, realConvID)
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

	if messages == nil {
		messages = make([]Message, 0)
	}

	return ConversationDetails{
		ConversationID: realConvID,
		Messages:       messages,
	}, nil
}
