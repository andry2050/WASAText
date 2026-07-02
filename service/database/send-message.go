package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) SendMessage(targetOrConvID string, senderID string, content string, isPhoto bool) (Message, error) {
	var realConvID string
	var isParticipant bool

	// 1. Controlla se targetOrConvID è un Gruppo o una Chat già avviata
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM conversations WHERE convid = ?)`, targetOrConvID).Scan(&isParticipant)

	if err == nil && isParticipant {
		// È una chat esistente. Verifichiamo se il mittente ne fa parte
		err = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM participants WHERE convid = ? AND userid = ?)`, targetOrConvID, senderID).Scan(&isParticipant)
		if err != nil || !isParticipant {
			return Message{}, fmt.Errorf("l'utente non fa parte di questa chat")
		}
		realConvID = targetOrConvID
	} else {
		// 2. Non è una chat esistente. Se l'ID appartiene a un UTENTE crea una chat diretta
		var targetExists bool
		err = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE userid = ?)`, targetOrConvID).Scan(&targetExists)
		if err != nil || !targetExists {
			return Message{}, fmt.Errorf("conversazione o utente inesistente")
		}

		// Genera un ID univoco per la chat (in ordine alfabetico per evitare duplicati)
		if senderID < targetOrConvID {
			realConvID = senderID + "_" + targetOrConvID
		} else {
			realConvID = targetOrConvID + "_" + senderID
		}

		// Crea la conversazione nel DB se non esiste già
		_, _ = db.c.Exec(`INSERT OR IGNORE INTO conversations (convid, type, name, photo_url) VALUES (?, 'direct', '', '')`, realConvID)
		_, _ = db.c.Exec(`INSERT OR IGNORE INTO participants (convid, userid) VALUES (?, ?)`, realConvID, senderID)
		_, _ = db.c.Exec(`INSERT OR IGNORE INTO participants (convid, userid) VALUES (?, ?)`, realConvID, targetOrConvID)
	}

	// Genera l'ID del messaggio e il timestamp
	msgUUID, err := uuid.NewV4()
	if err != nil {
		return Message{}, fmt.Errorf("errore generazione id messaggio: %w", err)
	}
	msgID := msgUUID.String()
	timestamp := time.Now().UTC()
	status := "sent"

	// Inserisce il messaggio
	query := `INSERT INTO messages (msgid, convid, senderid, content, is_photo, status, timestamp) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.c.Exec(query, msgID, realConvID, senderID, content, isPhoto, status, timestamp)
	if err != nil {
		return Message{}, fmt.Errorf("errore inserimento messaggio: %w", err)
	}

	// Recupera le informazioni del mittente
	var sender User
	var photoURL sql.NullString
	err = db.c.QueryRow(`SELECT username, photo_url FROM users WHERE userid = ?`, senderID).Scan(&sender.Username, &photoURL)
	if err != nil {
		return Message{}, fmt.Errorf("errore lettura dati mittente: %w", err)
	}
	sender.UserID = senderID
	if photoURL.Valid {
		sender.PhotoURL = photoURL.String
	}

	msg := Message{
		MessageID:      msgID,
		ConversationID: realConvID,
		SenderID:       senderID,
		Sender:         sender,
		Content:        content,
		IsPhoto:        isPhoto,
		Status:         status,
		Timestamp:      timestamp,
		Reactions:      make([]Reaction, 0),
	}

	return msg, nil
}
