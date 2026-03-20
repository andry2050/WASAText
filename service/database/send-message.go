package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// SendMessage salva un nuovo messaggio nel database e lo restituisce formattato
func (db *appdbimpl) SendMessage(convID string, senderID string, content string, isPhoto bool) (Message, error) {
	// Controlla se l'utente fa parte della conversazione
	var isParticipant bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM participants WHERE convid = ? AND userid = ?)`, convID, senderID).Scan(&isParticipant)
	if err != nil || !isParticipant {
		return Message{}, fmt.Errorf("l'utente non fa parte di questa chat")
	}

	// Genera l'ID del messaggio e il timestamp
	msgUUID, err := uuid.NewV4()
	if err != nil {
		return Message{}, fmt.Errorf("errore generazione id messaggio: %w", err)
	}
	msgID := msgUUID.String()
	timestamp := time.Now().UTC()

	// Imposta lo stato iniziale (una spunta = sent/received)
	status := "sent"

	// Inserisce il messaggio nella tabella
	query := `INSERT INTO messages (msgid, convid, senderid, content, is_photo, status, timestamp) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.c.Exec(query, msgID, convID, senderID, content, isPhoto, status, timestamp)
	if err != nil {
		return Message{}, fmt.Errorf("errore inserimento messaggio: %w", err)
	}

	// Recupera le informazioni dell'utente per costruire la risposta completa
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
		ConversationID: convID,
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
