package api

import (
	"time"

	"github.com/andry2050/WASAText/service/database"
)

type Photo struct {
	PhotoID   int       `json:"photoid"`
	UserID    int       `json:"userid"`
	User      User      `json:"user"`
	Image     string    `json:"image"`
	Timestamp time.Time `json:"timestamp"`
}

func (p *Photo) FromDatabase(photo database.Photo) {
	p.PhotoID = photo.PhotoID
	p.UserID = photo.UserID
	p.Image = photo.Image
	p.Timestamp = photo.Timestamp
}

func (p *Photo) ToDatabase() database.Photo {
	return database.Photo{
		PhotoID:   p.PhotoID,
		UserID:    p.UserID,
		Image:     p.Image,
		Timestamp: p.Timestamp,
	}
}
