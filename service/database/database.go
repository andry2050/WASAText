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
<<<<<<< HEAD

=======
>>>>>>> 226708c2193c4ffc194ad5e11414bb2dfcf65d82
package database

import (
	"database/sql"
	"errors"
	"fmt"
<<<<<<< HEAD
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
	ConversationID string `json:"id"`
	Type           string `json:"type"` // "direct" o "group"
	Name           string `json:"name,omitempty"`
	PhotoURL       string `json:"photo_url,omitempty"`
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
	IsPhoto        bool       `json:"is_photo"`
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

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// User operations
	DoLogin(username string) (string, error)
	SetMyUserName(userID string, newName string) error
	SetMyPhoto(userID string, photoPath string) error
	SearchUsers(query string) ([]User, error)
	
	// Conversation operations
	GetMyConversations(userID string) ([]Conversation, error)
	GetConversation(conversationID string) (Conversation, error)
	
	// Message operations
	SendMessage(msg Message) (Message, error)
	ForwardMessage(msgID string, targetConvID string, senderID string) (Message, error)
	DeleteMessage(msgID string, userID string) error
	
	// Reactions (Comments)
	CommentMessage(reaction Reaction) (Reaction, error)
	UncommentMessage(reactionID string, userID string) error

	// Group operations
	CreateGroup(group Conversation, memberIDs []string) (Conversation, error)
	AddToGroup(groupID string, userID string) error
	LeaveGroup(groupID string, userID string) error
	SetGroupName(groupID string, name string) error
	SetGroupPhoto(groupID string, photoPath string) error
=======
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
>>>>>>> 226708c2193c4ffc194ad5e11414bb2dfcf65d82

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
<<<<<<< HEAD
func New(db *sql.DB) (AppDatabase, error) {

=======
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
>>>>>>> 226708c2193c4ffc194ad5e11414bb2dfcf65d82
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

<<<<<<< HEAD
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
		content TEXT NOT NULL,
		is_photo BOOLEAN NOT NULL DEFAULT 0,
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
=======
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
>>>>>>> 226708c2193c4ffc194ad5e11414bb2dfcf65d82
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
