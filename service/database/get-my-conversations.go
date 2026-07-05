package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]Conversation, error) {
	// Usiamo delle sotto-query COALESCE. Se non ci sono messaggi, restituirà una stringa vuota
	// anziché nascondere l'intera conversazione! Inoltre rimaniamo agganciati ai partecipanti.
	query := `
		SELECT 
			c.convid, 
			c.type, 
			COALESCE(c.name, ''), 
			COALESCE(c.photo_url, ''),
			COALESCE((SELECT content FROM messages WHERE convid = c.convid ORDER BY timestamp DESC LIMIT 1), ''),
			COALESCE((SELECT timestamp FROM messages WHERE convid = c.convid ORDER BY timestamp DESC LIMIT 1), '')
		FROM conversations c
		INNER JOIN participants p ON c.convid = p.convid
		WHERE p.userid = ?
	`
	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("errore recupero conversazioni dal db: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		err = rows.Scan(&c.ConversationID, &c.Type, &c.Name, &c.PhotoURL, &c.LastMessagePreview, &c.LastMessageTimestamp)
		if err != nil {
			return nil, fmt.Errorf("errore scan conversazione: %w", err)
		}

		// ✨ MAGIA PER LE CHAT DIRETTE:
		// Se è una chat diretta, il campo 'name' nel DB è vuoto.
		// Dobbiamo estrarre il nome e la foto dell'ALTRA persona presente nella chat!
		if c.Type == "direct" {
			var otherUsername, otherPhoto sql.NullString
			errOther := db.c.QueryRow(`
				SELECT u.username, u.photo_url 
				FROM participants p
				INNER JOIN users u ON p.userid = u.userid
				WHERE p.convid = ? AND p.userid != ?
			`, c.ConversationID, userID).Scan(&otherUsername, &otherPhoto)

			if errOther == nil {
				if otherUsername.Valid {
					c.Name = otherUsername.String
				}
				if otherPhoto.Valid {
					c.PhotoURL = otherPhoto.String
				}
			}
		}

		conversations = append(conversations, c)
	}

	// Evitiamo di restituire 'null' al frontend se l'array è vuoto, restituiamo invece []
	if conversations == nil {
		conversations = make([]Conversation, 0)
	}

	return conversations, nil
}
