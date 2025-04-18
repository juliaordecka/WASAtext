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
	"log"
	"time"
)

var ErrUserDoesNotExist = errors.New("User does not exist")

type User struct {
	Id           uint64 `json:"id"`
	Username     string `json:"username"`
	ProfilePhoto string `json:"profilePhoto,omitempty"`
}

type Message struct {
	MessageId      int       `json:"messageId"`
	Text           string    `json:"text"`
	SendTime       time.Time `json:"sendTime"`
	Status         string    `json:"status"`
	SenderId       uint64    `json:"senderId"`
	RecipientId    uint64    `json:"recipientId"`
	ConversationId int       `json:"conversationId"`
	Photo          string    `json:"photo"`
}

type Conversation struct {
	ConversationId int    `json:"conversationId"`
	GroupId        int    `json:"GroupId"`
	LastMessageId  int    `json:"lastMessageId"`
	Name           string `json:"name"`
}

// New structs

type ConversationPreview struct {
	ConversationId  int       `json:"conversationId"`
	Name            string    `json:"name"`
	Photo           string    `json:"photo,omitempty"`
	LastMessageTime time.Time `json:"lastMessageTime"`
	LastMessageText string    `json:"lastMessageText"`
	IsPhoto         bool      `json:"isPhoto"`
	IsGroup         bool      `json:"isGroup"`
}

type ConversationDetails struct {
	ConversationId int                   `json:"conversationId"`
	Name           string                `json:"name"`
	Photo          string                `json:"photo,omitempty"`
	IsGroup        bool                  `json:"isGroup"`
	Messages       []MessageWithComments `json:"messages"`
}

type MessageWithComments struct {
	Message
	SenderUsername string    `json:"senderUsername"`
	Comments       []Comment `json:"comments"`
}

type Comment struct {
	UserId   uint64 `json:"userId"`
	Username string `json:"username"`
	Emoji    string `json:"emoji"`
}

// End of new structs

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	CreateUser(User) (User, error)
	SetUsername(User, string) (User, error)
	GetOrCreateDirectConversation(userId, recipientId uint64) (int, error)
	CreateMessage(Message) (Message, error)
	CheckIfConversationExists(conversationId int) (bool, error)
	CreateConversation(userId uint64, conversationId int) (Conversation, error)
	UpdateLastMessage(messageId int, conversationId int) error
	// The group implementation
	GetUsernameById(userId uint64) (string, error)
	CreateGroup(name string, creatorId uint64) (Conversation, error)
	AddUserToGroup(username string, groupId int) error
	GetUserIdByUsername(username string) (uint64, error)
	IsUserInGroup(userId uint64, groupId int) (bool, error)
	DeleteGroup(groupId int) error
	// Messaging to groups
	GetConversationIdByName(name string) (int, error)
	GetRecipientIdByUsername(username string) (uint64, error)
	// Leave group and set group name
	LeaveGroup(userId uint64, groupId int) error
	SetGroupName(groupId int, newName string) error
	// Comments
	ForwardMessage(messageId int, userId uint64, targetConvId int) (Message, error)
	DeleteMessage(messageId int, userId uint64) error
	CommentMessage(messageId int, userId uint64, emoji string) error
	UncommentMessage(messageId int, userId uint64) error
	IsMessageOwner(messageId int, userId uint64) (bool, error)
	// Last functions
	SetUserPhoto(userId uint64, photoData string) error
	SetGroupPhoto(groupId int, photoData string) error
	GetConversations(userId uint64) ([]ConversationPreview, error)
	GetConversationDetails(convId int, userId uint64) (ConversationDetails, error)
	SearchUsers(query string) ([]User, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Creating 'users' table...")
		usersDatabase := `CREATE TABLE users (
			Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			Username TEXT NOT NULL UNIQUE,
			ProfilePhoto TEXT
			);`
		_, err = db.Exec(usersDatabase)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		} else {
			log.Println("'users' table successfully created.")
		}
	} else if err != nil {
		return nil, fmt.Errorf("error checking table existence: %w", err)
	} else {
		log.Println("'users' table already exists.")
	}

	// Message table
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Creating 'messages' table...")
		messagesDatabase := `CREATE TABLE messages (
            MessageId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            ConversationId INTEGER NOT NULL,
            Text TEXT NOT NULL,
            SendTime DATETIME NOT NULL,
            Status TEXT NOT NULL,
            SenderId INTEGER NOT NULL,
            RecipientId INTEGER NOT NULL,
            Photo TEXT,
            FOREIGN KEY (ConversationId) REFERENCES conversations(ConversationId)
        );`
		_, err = db.Exec(messagesDatabase)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		} else {
			log.Println("'messages' table successfully created.")
		}
	} else if err != nil {
		return nil, fmt.Errorf("error checking table existence: %w", err)
	} else {
		log.Println("'messages' table already exists.")
	}

	// Conversation table
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversations';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Creating 'conversations' table...")
		conversationsDatabase := `CREATE TABLE conversations (
			ConversationId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			GroupId INTEGER NOT NULL,
			LastMessageId INTEGER,
			Name TEXT,
			GroupPhoto TEXT,
			FOREIGN KEY (LastMessageId) REFERENCES messages(MessageId)
		);`

		_, err = db.Exec(conversationsDatabase)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		} else {
			log.Println("'conversations' table successfully created.")
		}
	} else if err != nil {
		return nil, fmt.Errorf("error checking table existence: %w", err)
	} else {
		log.Println("'conversations' table already exists.")
	}

	// Participants table
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='participants';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Creating 'participants' table...")
		participantsDatabase := `CREATE TABLE participants (
            ConversationId INTEGER NOT NULL,
            UserId INTEGER NOT NULL,
            PRIMARY KEY (ConversationId, UserId),
            FOREIGN KEY (ConversationId) REFERENCES conversations(ConversationId),
            FOREIGN KEY (UserId) REFERENCES users(Id)
        );`
		_, err = db.Exec(participantsDatabase)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		} else {
			log.Println("'participants' table successfully created.")
		}
	} else if err != nil {
		return nil, fmt.Errorf("error checking table existence: %w", err)
	} else {
		log.Println("'participants' table already exists.")
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Creating 'comments' table...")
		commentsDatabase := `CREATE TABLE comments (
        CommentId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        MessageId INTEGER NOT NULL,
        UserId INTEGER NOT NULL,
        Emoji TEXT NOT NULL,
        FOREIGN KEY (MessageId) REFERENCES messages(MessageId),
        FOREIGN KEY (UserId) REFERENCES users(Id),
        UNIQUE(MessageId, UserId)
    );`
		_, err = db.Exec(commentsDatabase)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		} else {
			log.Println("'comments' table successfully created.")
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
