package api

import (
	"regexp"

	"github.com/andry2050/wasa/service/database"
)

type User struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
}

func (u *User) FromDatabase(user database.User) {
	u.UserID = user.UserID
	u.Username = user.Username

}

func (u *User) ToDatabase() database.User {
	return database.User{
		UserID:   u.UserID,
		Username: u.Username,
	}
}

func (u *User) IsValid() bool {
	username := u.Username
	validUser := regexp.MustCompile(`^.*?$`)
	return validUser.MatchString(username)
}