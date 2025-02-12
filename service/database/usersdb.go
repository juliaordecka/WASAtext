package database

import (
    "fmt"
)

func (db *appdbimpl) SetUserPhoto(userId uint64, photoData string) error {
    _, err := db.c.Exec("UPDATE users SET ProfilePhoto = ? WHERE Id = ?", 
        photoData, userId)
    return err
}

func (db *appdbimpl) SetGroupPhoto(groupId int, photoData string) error {
    _, err := db.c.Exec("UPDATE conversations SET GroupPhoto = ? WHERE ConversationId = ? AND GroupId = 1", 
        photoData, groupId)
    return err
}

func (db *appdbimpl) GetConversations(userId uint64) ([]ConversationPreview, error) {
    rows, err := db.c.Query(`
        SELECT 
            c.ConversationId,
            CASE 
                WHEN c.GroupId = 1 THEN c.Name
                ELSE u.Username
            END as Name,
            CASE 
                WHEN c.GroupId = 1 THEN c.GroupPhoto
                ELSE u.ProfilePhoto
            END as Photo,
            m.SendTime,
            m.Text,
            m.Photo IS NOT NULL as IsPhoto,
            c.GroupId = 1 as IsGroup
        FROM conversations c
        JOIN participants p ON c.ConversationId = p.ConversationId
        LEFT JOIN messages m ON c.LastMessageId = m.MessageId
        LEFT JOIN users u ON (c.GroupId = 0 AND p.UserId != ? AND u.Id = p.UserId)
        WHERE p.UserId = ?
        ORDER BY m.SendTime DESC`, userId, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var conversations []ConversationPreview
    for rows.Next() {
        var conv ConversationPreview
        err := rows.Scan(
            &conv.ConversationId,
            &conv.Name,
            &conv.Photo,
            &conv.LastMessageTime,
            &conv.LastMessageText,
            &conv.IsPhoto,
            &conv.IsGroup,
        )
        if err != nil {
            return nil, err
        }
        conversations = append(conversations, conv)
    }
    return conversations, nil
}

func (db *appdbimpl) GetConversationDetails(convId int, userId uint64) (ConversationDetails, error) {
    var conv ConversationDetails

    // Get conversation info
    err := db.c.QueryRow(`
        SELECT 
            c.ConversationId,
            CASE 
                WHEN c.GroupId = 1 THEN c.Name
                ELSE u.Username
            END as Name,
            CASE 
                WHEN c.GroupId = 1 THEN c.GroupPhoto
                ELSE u.ProfilePhoto
            END as Photo,
            c.GroupId = 1 as IsGroup
        FROM conversations c
        LEFT JOIN participants p ON c.ConversationId = p.ConversationId AND p.UserId != ?
        LEFT JOIN users u ON c.GroupId = 0 AND u.Id = p.UserId
        WHERE c.ConversationId = ?`, userId, convId).Scan(
        &conv.ConversationId,
        &conv.Name,
        &conv.Photo,
        &conv.IsGroup,
    )
    if err != nil {
        return conv, err
    }

    // Get messages with comments
    rows, err := db.c.Query(`
        SELECT 
            m.MessageId,
            m.Text,
            m.SendTime,
            m.Status,
            m.SenderId,
            m.RecipientId,
            m.Photo,
            u.Username as SenderUsername,
            cm.UserId as CommentUserId,
            cu.Username as CommentUsername,
            cm.Emoji
        FROM messages m
        JOIN users u ON m.SenderId = u.Id
        LEFT JOIN comments cm ON m.MessageId = cm.MessageId
        LEFT JOIN users cu ON cm.UserId = cu.Id
        WHERE m.ConversationId = ?
        ORDER BY m.SendTime DESC`, convId)
    if err != nil {
        return conv, err
    }
    defer rows.Close()

    messageMap := make(map[int]*MessageWithComments)
    for rows.Next() {
        var msg MessageWithComments
        var commentUserId *uint64
        var commentUsername, emoji *string

        err := rows.Scan(
            &msg.MessageId,
            &msg.Text,
            &msg.SendTime,
            &msg.Status,
            &msg.SenderId,
            &msg.RecipientId,
            &msg.Photo,
            &msg.SenderUsername,
            &commentUserId,
            &commentUsername,
            &emoji,
        )
        if err != nil {
            return conv, err
        }

        existing, exists := messageMap[msg.MessageId]
        if !exists {
            messageMap[msg.MessageId] = &msg
            existing = &msg
        }

        if commentUserId != nil {
            existing.Comments = append(existing.Comments, Comment{
                UserId:   *commentUserId,
                Username: *commentUsername,
                Emoji:   *emoji,
            })
        }
    }

    for _, msg := range messageMap {
        conv.Messages = append(conv.Messages, *msg)
    }

    return conv, nil
}

func (db *appdbimpl) SearchUsers(query string) ([]User, error) {
    rows, err := db.c.Query(`
        SELECT Id, Username, ProfilePhoto 
        FROM users 
        WHERE Username LIKE ?
        ORDER BY Username`, 
        fmt.Sprintf("%%%s%%", query))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.Id, &user.Username, &user.ProfilePhoto)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

