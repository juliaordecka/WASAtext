package api

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"time"
)

// user struct

type User struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}

func (u *User) FromDatabase(user database.User) {
	u.Id = user.Id
	u.Username = user.Username
}

func (u *User) ToDatabase() database.User {
	return database.User{
		Id:       u.Id,
		Username: u.Username,
	}
}

// Message struct

type Message struct {
	MessageId      int       `json:"messageId"`
	ConversationId int       `json:"conversationId"`
	Text           string    `json:"text,omitempty"`
	SendTime       time.Time `json:"sendTime"`
	Status         string    `json:"status"`
	SenderId       uint64    `json:"senderId"`
	RecipientId    uint64    `json:"recipientId,omitempty"`
	RecipientUsername string `json:"recipientUsername,omitempty"`
    ConversationName string `json:"conversationName,omitempty"`
	Photo          string    `json:"photo,omitempty"`
}

func (m *Message) FromDatabase(dbMsg database.Message) {
	m.MessageId = dbMsg.MessageId
	m.ConversationId = dbMsg.ConversationId
	m.Text = dbMsg.Text
	m.SendTime = dbMsg.SendTime
	m.Status = dbMsg.Status
	m.SenderId = dbMsg.SenderId // Convert to uint64 from int
	m.RecipientId = dbMsg.RecipientId
	m.Photo = dbMsg.Photo
}

// ToDatabase converts an api Message into a database Message
func (m *Message) ToDatabase() database.Message {
	return database.Message{
		MessageId:      m.MessageId,
		ConversationId: m.ConversationId,
		Text:           m.Text,
		SendTime:       m.SendTime,
		Status:         m.Status,
		SenderId:       m.SenderId, // Convert to int from uint64
		RecipientId:    m.RecipientId,
		Photo:          m.Photo,
	}
}

// Group struct
type Group struct {
	GroupId int    `json:"groupId"`
	Name    string `json:"name"`
}

// Conversation struct

type Conversation struct {
	ConversationId int    `json:"conversationId"`
	GroupId        int    `json:"GroupId"`
	LastMessageId  int    `json:"lastMessageId"`
	Name           string `json:"name,omitempty"`
	Participants   []User `json:"participants,omitempty"`
	IsGroup        bool   `json:"isGroup"`
}
