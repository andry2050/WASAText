/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

// Modelli Dati per WASAText

type User struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	PhotoURL string `json:"photo_url"`
}

type Conversation struct {
	ConversationID       string `json:"id"`
	Type                 string `json:"type"` // "direct" o "group"
	Name                 string `json:"name"`
	PhotoURL             string `json:"photo_url"`
	LastMessagePreview   string `json:"last_message_preview"`
	LastMessageTimestamp string `json:"last_message_timestamp"`
}

type ConversationDetails struct {
	ConversationID string    `json:"id"`
	Messages       []Message `json:"messages"`
}

type Participant struct {
	ConversationID string
	UserID         string
}

type Message struct {
	MessageID      string     `json:"id"`
	ConversationID string     `json:"-"`
	SenderID       string     `json:"-"`
	Sender         User       `json:"sender"`
	Content        string     `json:"content"`
	PhotoURL       string     `json:"photo_url"`
	Status         string     `json:"status"` // "sent", "received", "read"
	Timestamp      time.Time  `json:"timestamp"`
	Reactions      []Reaction `json:"reactions"`
}

type Reaction struct {
	ReactionID string `json:"id"`
	MessageID  string `json:"-"`
	UserID     string `json:"-"`
	User       User   `json:"user"`
	Emoji      string `json:"emoji"`
}

type Group struct {
	GroupID  string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
	Members  []User `json:"members"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// User operations
	DoLogin(username string) (string, error)
	SetMyUserName(userID string, newName string) error
	SetMyPhoto(userID string, photoPath string) error
	SearchUsers(query string) ([]User, error)

	// Conversation operations
	GetMyConversations(userID string) ([]Conversation, error)
	GetConversation(conversationID string, userID string) (ConversationDetails, error)

	// Message operations
	SendMessage(convID string, senderID string, content string, photoURL string) (Message, error)
	ForwardMessage(msgID string, targetConvID string, senderID string) (Message, error)
	DeleteMessage(msgID string, userID string) error
	MarkMessagesAsRead(convID string, myID string) error

	// Reactions (Comments)
	CommentMessage(messageID string, userID string, reqBodyEmoji string) (Reaction, error)
	UncommentMessage(messageID string, commentID string, userID string) error

	// Group operations
	CreateGroup(groupName string, memberIDs []string, creatorID string) (Group, error)
	AddToGroup(groupID string, userID string, requesterID string) error
	LeaveGroup(groupID string, userID string) error
	SetGroupName(groupID string, name string, requesterID string) error
	SetGroupPhoto(groupID string, photoPath string, requesterID string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
func New(db *sql.DB) (AppDatabase, error) {

	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var err error

	_, errPramga := db.Exec(`PRAGMA foreign_keys = ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// USER TABLE
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		userid TEXT NOT NULL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		photo_url TEXT
		);`)
	if err != nil {
		return nil, err
	}

	// CONVERSATIONS TABLE (usata per chat dirette e gruppi)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS conversations (
		convid TEXT NOT NULL PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT,
		photo_url TEXT
		);`)
	if err != nil {
		return nil, err
	}

	// CONVERSATION PARTICIPANTS TABLE
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS participants (
		convid TEXT NOT NULL,
		userid TEXT NOT NULL,
		PRIMARY KEY (convid, userid),
		FOREIGN KEY (convid) REFERENCES conversations(convid) ON DELETE CASCADE,
		FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
		);`)
	if err != nil {
		return nil, err
	}

	// MESSAGES TABLE
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		msgid TEXT NOT NULL PRIMARY KEY,
		convid TEXT NOT NULL,
		senderid TEXT NOT NULL,
		content TEXT,          
		photo_url TEXT,        
		status TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		FOREIGN KEY (convid) REFERENCES conversations(convid) ON DELETE CASCADE,
		FOREIGN KEY (senderid) REFERENCES users(userid) ON DELETE CASCADE
		);`)
	if err != nil {
		return nil, err
	}

	// REACTIONS TABLE
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reactions (
		reactionid TEXT NOT NULL PRIMARY KEY,
		msgid TEXT NOT NULL,
		userid TEXT NOT NULL,
		emoji TEXT NOT NULL,
		FOREIGN KEY (msgid) REFERENCES messages(msgid) ON DELETE CASCADE,
		FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
		);`)
	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
