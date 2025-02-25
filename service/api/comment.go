package api

import (
	"github.com/andry2050/wasa/service/database"
)

type Comment struct {
	CommentID int    `json:"commentid"`
	UserID    int    `json:"userid"`
	User      User   `json:"user"`
	PostID    int    `json:"postid"`
	Text      string `json:"text"`
}

func (c *Comment) FromDatabase(comment database.Comment) {
	c.CommentID = comment.CommentID
	c.UserID = comment.UserID
	c.PostID = comment.PostID
	c.Text = comment.Text

}

func (c *Comment) ToDatabase() database.Comment {
	return database.Comment{
		CommentID: c.CommentID,
		UserID:    c.UserID,
		PostID:    c.PostID,
		Text:      c.Text,
	}
}
