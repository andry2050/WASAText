package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) SendMessage(targetOrConvID string, senderID string, content string, isPhoto bool) (Message, error) {
	var realConvID string
	var isConv bool

	// 1. Controlliamo se targetOrConvID è una chat già avviata
	_ = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM conversations WHERE convid = ?)`, targetOrConvID).Scan(&isConv)

	if isConv {
		var isParticipant bool
		_ = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM participants WHERE convid = ? AND userid = ?)`, targetOrConvID, senderID).Scan(&isParticipant)
		if !isParticipant {
			return Message{}, fmt.Errorf("accesso negato a questa chat")
		}
		realConvID = targetOrConvID
	} else {
		// 2. Creazione manuale e sicura della chat diretta
		var targetExists bool
		err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE userid = ?)`, targetOrConvID).Scan(&targetExists)
		if err != nil || !targetExists {
			return Message{}, fmt.Errorf("utente inesistente")
		}

		if senderID < targetOrConvID {
			realConvID = senderID + "_" + targetOrConvID
		} else {
			realConvID = targetOrConvID + "_" + senderID
		}

		// Verifichiamo se l'avevamo già creata per evitare errori di collisione
		var chatExists bool
		_ = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM conversations WHERE convid = ?)`, realConvID).Scan(&chatExists)

		if !chatExists {
			_, _ = db.c.Exec(`INSERT INTO conversations (convid, type, name, photo_url) VALUES (?, 'direct', '', '')`, realConvID)
			_, _ = db.c.Exec(`INSERT INTO participants (convid, userid) VALUES (?, ?)`, realConvID, senderID)
			_, _ = db.c.Exec(`INSERT INTO participants (convid, userid) VALUES (?, ?)`, realConvID, targetOrConvID)
		}
	}

	// 3. Generazione e Inserimento del messaggio
	msgUUID, err := uuid.NewV4()
	if err != nil {
		return Message{}, fmt.Errorf("errore id: %w", err)
	}
	msgID := msgUUID.String()
	timestamp := time.Now().UTC()
	status := "sent"

	query := `INSERT INTO messages (msgid, convid, senderid, content, is_photo, status, timestamp) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.c.Exec(query, msgID, realConvID, senderID, content, isPhoto, status, timestamp)
	if err != nil {
		return Message{}, fmt.Errorf("errore inserimento messaggio: %w", err)
	}

	// 4. Recupero dati utente mittente
	var sender User
	var photoURL sql.NullString
	_ = db.c.QueryRow(`SELECT username, photo_url FROM users WHERE userid = ?`, senderID).Scan(&sender.Username, &photoURL)
	sender.UserID = senderID
	if photoURL.Valid {
		sender.PhotoURL = photoURL.String
	}

	return Message{
		MessageID:      msgID,
		ConversationID: realConvID,
		SenderID:       senderID,
		Sender:         sender,
		Content:        content,
		IsPhoto:        isPhoto,
		Status:         status,
		Timestamp:      timestamp,
		Reactions:      make([]Reaction, 0),
	}, nil
}
