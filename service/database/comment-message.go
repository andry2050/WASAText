package database

func (db *appdbimpl) CommentMessage(c Comment) error {
	res, err := db.c.Exec(`INSERT INTO comments (commentid, userid, messageid, text) VALUES (NULL, ?, ?, ?)`, c.CommentID, c.UserID, c.MessageID, c.Text)
	if err != nil {
		return err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	c.CommentID = int(lastInsertID)
	return nil
}

func (db *appdbimpl) UncommentMessage(c Comment) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE commentid = ?`, c.CommentID)
	if err != nil {
		return err
	}
	return nil
}
