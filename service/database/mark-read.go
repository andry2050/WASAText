package database

func (db *appdbimpl) MarkMessagesAsRead(convID string, userID string) error {
	_, err := db.c.Exec(`
		UPDATE messages 
		SET status = 'read' 
		WHERE conversation_id = ? AND sender_id != ?`,
		convID, userID,
	)
	return err
}
