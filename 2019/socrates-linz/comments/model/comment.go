package model

import (
	"database/sql"
	"fmt"
	"time"
)

// Comment holds a comment.
type Comment struct {
	Mail    string
	Created time.Time
	Message string
}

func (p Comment) String() string {
	return fmt.Sprintf("{Mail: %s, Created: %s, Message: %s}", p.Mail, p.Created, p.Message)
}

// CommentInit initializes the comments table.
func CommentInit(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE comments(mail TEXT, created datetime, message TEXT)")

	return err
}

// CommentAdd adds a comment to the database.
func CommentAdd(db *sql.DB, mail string, message string) error {
	_, err := db.Exec("INSERT INTO comments(mail, created, message) VALUES('" + mail + "', datetime('now'), '" + message + "')")

	return err
}

// CommentAll returns all comments ordered in by their creation time.
func CommentAll(db *sql.DB) ([]*Comment, error) {
	rows, err := db.Query("SELECT mail, created, message FROM comments ORDER BY created DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Mail, &comment.Created, &comment.Message)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}
