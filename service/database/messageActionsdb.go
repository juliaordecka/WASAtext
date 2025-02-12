package database

import(
	"time"
)

func (db *appdbimpl) ForwardMessage(messageId int, userId uint64, targetConvId int) (Message, error) {
    // Get original message
    var msg Message
    err := db.c.QueryRow(`
        SELECT Text, Status, SenderId, Photo 
        FROM messages 
        WHERE MessageId = ?`, messageId).Scan(&msg.Text, &msg.Status, &msg.SenderId, &msg.Photo)
    if err != nil {
        return Message{}, err
    }

    // Create new message
    msg.ConversationId = targetConvId
    msg.SenderId = userId
    msg.SendTime = time.Now()
    
    return db.CreateMessage(msg)
}

func (db *appdbimpl) DeleteMessage(messageId int, userId uint64) error {
    // Start transaction
    tx, err := db.c.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Delete comments first
    _, err = tx.Exec("DELETE FROM comments WHERE MessageId = ?", messageId)
    if err != nil {
        return err
    }

    // Delete message
    _, err = tx.Exec("DELETE FROM messages WHERE MessageId = ? AND SenderId = ?", 
        messageId, userId)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (db *appdbimpl) CommentMessage(messageId int, userId uint64, emoji string) error {
    _, err := db.c.Exec(`
        INSERT INTO comments (MessageId, UserId, Emoji) 
        VALUES (?, ?, ?)`, messageId, userId, emoji)
    return err
}

func (db *appdbimpl) UncommentMessage(messageId int, userId uint64) error {
    _, err := db.c.Exec(`
        DELETE FROM comments 
        WHERE MessageId = ? AND UserId = ?`, messageId, userId)
    return err
}

func (db *appdbimpl) IsMessageOwner(messageId int, userId uint64) (bool, error) {
    var exists bool
    err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM messages 
            WHERE MessageId = ? AND SenderId = ?
        )`, messageId, userId).Scan(&exists)
    return exists, err
}

