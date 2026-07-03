package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]Conversation, error) {
	// Il comando NULLIF è fondamentale: impedisce che i nomi siano stringhe vuote (rendendoli invisibili)
	query := `
		SELECT
			c.convid,
			c.type,
			COALESCE(NULLIF(c.name, ''), u.username, 'Utente Sconosciuto') AS name,
			COALESCE(NULLIF(c.photo_url, ''), u.photo_url, '') AS photo_url,
			lm.content,
			lm.is_photo,
			lm.timestamp
		FROM conversations c
		INNER JOIN participants p1 ON c.convid = p1.convid AND p1.userid = ?
		LEFT JOIN participants p2 ON c.convid = p2.convid AND p2.userid != ? AND c.type = 'direct'
		LEFT JOIN users u ON p2.userid = u.userid
		LEFT JOIN (
			SELECT m1.convid, m1.content, m1.is_photo, m1.timestamp
			FROM messages m1
			INNER JOIN (
				SELECT convid, MAX(timestamp) as max_ts FROM messages GROUP BY convid
			) m2 ON m1.convid = m2.convid AND m1.timestamp = m2.max_ts
			GROUP BY m1.convid
		) lm ON c.convid = lm.convid
		ORDER BY lm.timestamp DESC
	`

	rows, err := db.c.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("errore query conversazioni: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation

	for rows.Next() {
		var c Conversation
		var msgContent sql.NullString
		var msgIsPhoto sql.NullBool
		var msgTimestamp sql.NullTime

		err = rows.Scan(&c.ConversationID, &c.Type, &c.Name, &c.PhotoURL, &msgContent, &msgIsPhoto, &msgTimestamp)
		if err != nil {
			return nil, fmt.Errorf("errore lettura riga: %w", err)
		}

		if msgTimestamp.Valid {
			c.LastMessageTimestamp = msgTimestamp.Time.Format(time.RFC3339)
			if msgIsPhoto.Valid && msgIsPhoto.Bool {
				c.LastMessagePreview = "📷 Foto"
			} else if msgContent.Valid {
				c.LastMessagePreview = msgContent.String
			}
		} else {
			c.LastMessageTimestamp = ""
			c.LastMessagePreview = "Nessun messaggio"
		}

		conversations = append(conversations, c)
	}

	if conversations == nil {
		conversations = make([]Conversation, 0)
	}
	return conversations, nil
}
