package database

import (
	"log"
)

func (db *appdbimpl) CreateMessage(m Message) (Message, error) {
	// Insert the message into the database using the correct SenderId
	log.Printf("Attempting to create message: %+v", m)
	res, err := db.c.Exec("INSERT INTO messages (ConversationId, SenderId, RecipientId, Text, Status, SendTime, Photo) VALUES (?, ?, ?, ?, ?, ?, ?)",
		m.ConversationId, m.SenderId, m.RecipientId, m.Text, m.Status, m.SendTime, m.Photo)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		return m, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return m, err
	}
	m.MessageId = int(lastInsertID)

	return m, nil
}
