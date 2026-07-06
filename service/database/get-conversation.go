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
		return ConversationDetails{
			ConversationID: realConvID,
			Messages:       make([]Message, 0),
		}, nil
	}

	// Quando un utente apre la chat, tutti i messaggi inviati dall'altra persona passano a read
	_, _ = db.c.Exec(`UPDATE messages SET status = 'read' WHERE convid = ? AND senderid != ? AND status != 'read'`, realConvID, userID)

	// Prende i messaggi
	query := `
		SELECT 
			m.msgid, m.senderid, u.username, u.photo_url, 
			m.content, m.photo_url, m.status, m.timestamp
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
		var senderPhotoURL sql.NullString
		var msgPhotoURL sql.NullString

		err = rows.Scan(
			&msg.MessageID, &msg.SenderID, &msg.Sender.Username, &senderPhotoURL,
			&msg.Content, &msgPhotoURL, &msg.Status, &msgTime,
		)
		if err != nil {
			return ConversationDetails{}, fmt.Errorf("errore scan messaggio: %w", err)
		}

		// Assegna i dati letti alla struttura del messaggio
		msg.Sender.UserID = msg.SenderID
		if senderPhotoURL.Valid {
			msg.Sender.PhotoURL = senderPhotoURL.String
		}

		// Salva il percorso della foto allegata al messaggio
		if msgPhotoURL.Valid {
			msg.PhotoURL = msgPhotoURL.String
		}

		msg.Timestamp = msgTime
		// Recupera le reazioni associate a questo specifico messaggio
		reacRows, errReac := db.c.Query(`
			SELECT r.reactionid, r.emoji, u.userid, u.username 
			FROM reactions r 
			INNER JOIN users u ON r.userid = u.userid 
			WHERE r.msgid = ?
		`, msg.MessageID)

		var reactions []Reaction
		if errReac == nil {
			for reacRows.Next() {
				var r Reaction
				var rUserID, rUserName string
				reacRows.Scan(&r.ReactionID, &r.Emoji, &rUserID, &rUserName)
				r.User = User{UserID: rUserID, Username: rUserName}
				reactions = append(reactions, r)
			}
			reacRows.Close()
		}

		if reactions == nil {
			reactions = make([]Reaction, 0)
		}
		msg.Reactions = reactions

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
