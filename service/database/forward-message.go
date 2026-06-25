package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// ForwardMessage clona un messaggio esistente in una nuova conversazione
func (db *appdbimpl) ForwardMessage(originalMsgID string, targetConvID string, senderID string) (Message, error) {
	// Controlla se l'utente fa parte della chat di destinazione
	var isTargetParticipant bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM participants WHERE convid = ? AND userid = ?)`, targetConvID, senderID).Scan(&isTargetParticipant)
	if err != nil || !isTargetParticipant {
		return Message{}, fmt.Errorf("Utente non fa parte della chat di destinazione")
	}

	// Recupera il contenuto del messaggio originale e verifica che l'utente faccia parte della
	// chat da cui proviene il messaggio originale
	var content string
	var isPhoto bool
	var sourceConvID string

	querySource := `
		SELECT m.content, m.is_photo, m.convid 
		FROM messages m
		INNER JOIN participants p ON m.convid = p.convid
		WHERE m.msgid = ? AND p.userid = ?
	`
	err = db.c.QueryRow(querySource, originalMsgID, senderID).Scan(&content, &isPhoto, &sourceConvID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Message{}, fmt.Errorf("messaggio originale non trovato o accesso negato")
		}
		return Message{}, fmt.Errorf("errore recupero messaggio originale: %w", err)
	}

	// Crea il nuovo messaggio
	newMsgUUID, _ := uuid.NewV4()
	newMsgID := newMsgUUID.String()
	timestamp := time.Now().UTC()
	status := "sent" // Anche l'inoltro parte con lo stato di "inviato"

	insertQuery := `INSERT INTO messages (msgid, convid, senderid, content, is_photo, status, timestamp) 
	                VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.c.Exec(insertQuery, newMsgID, targetConvID, senderID, content, isPhoto, status, timestamp)
	if err != nil {
		return Message{}, fmt.Errorf("errore inserimento messaggio inoltrato: %w", err)
	}

	// Recupera le info dell'utente mittente per la risposta JSON
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

	newMsg := Message{
		MessageID:      newMsgID,
		ConversationID: targetConvID,
		SenderID:       senderID,
		Sender:         sender,
		Content:        content,
		IsPhoto:        isPhoto,
		Status:         status,
		Timestamp:      timestamp,
		Reactions:      make([]Reaction, 0),
	}

	return newMsg, nil
}
