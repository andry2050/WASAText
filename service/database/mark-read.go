package database

func (db *appdbimpl) MarkMessagesAsRead(targetOrConvID string, userID string) error {
	var realConvID string
	var isConv bool

	// Controlla se la conversazione esiste già
	_ = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM conversations WHERE convid = ?)`, targetOrConvID).Scan(&isConv)

	if isConv {
		realConvID = targetOrConvID
	} else {
		// Se non esiste, calcola l'ID univoco alfabetico (es. Marco_Matteo)
		if userID < targetOrConvID {
			realConvID = userID + "_" + targetOrConvID
		} else {
			realConvID = targetOrConvID + "_" + userID
		}
	}

	// Ora aggiorna i messaggi senza andare in errore 500, anche se la chat è vuota!
	_, err := db.c.Exec(`UPDATE messages SET status = 'read' WHERE convid = ? AND senderid != ? AND status != 'read'`, realConvID, userID)
	return err
}
