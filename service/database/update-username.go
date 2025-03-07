package database

func (db *appdbimpl) UpdateUsername(u User) error {
	res, err := db.c.Exec(`UPDATE users SET username=? WHERE userid=?`, u.Username, u.UserID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrUserDoesNotExist
	}
	return nil
}
