package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/3XBAT/coments-system/graph/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (r *CommentPostgres) Create(ctx context.Context, input *model.NewComment) (string, error) {

	newCommentId := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)

	logrus.Printf("ParentId - %s", input.ParentID)

	createPostQuery := fmt.Sprintf(`INSERT INTO %s (comment_id, post_id, parent_id, content, username, created_at)
	 VALUES ($1, $2, $3, $4, $5, $6)`, commentsTable)

	_, err := r.db.Exec(createPostQuery, newCommentId, input.PostID, input.ParentID, input.Content, input.Username, createdAt)

	if err != nil {
		logrus.Errorf("failed insert comment into db:%s", err.Error())
		return "", err
	}

	return newCommentId, nil
}

func (r *CommentPostgres) GetById(ctx context.Context, id string) (*model.Comment, error) {

	query := (`SELECT
		c.comment_id,
		c.post_id,
		c.parent_id,
		c.content,
		c.username,
		c.created_at,
		c.updated_at
	FROM
		comments c
	WHERE
		c.post_id = $1
	ORDER BY
		 c.created_at, c.parent_id;`)

	rows, err := r.db.Queryx(query, id)
	if err != nil {
		logrus.Errorf("failed receive data from database:%s", err.Error())
		return nil, errors.New("there are no comments with this id")
	}

	defer rows.Close()

	var comment model.Comment

	for rows.Next() {
		var updatedAt sql.NullString
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.Username, &comment.CreatedAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		comment.Replies = append(comment.Replies, comment.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentPostgres) GetPart(ctx context.Context, postID string, parentID string, limit *int, offset *int) ([]*model.Comment, error) {
	var comments []*model.Comment

	if *limit <= 0 || *offset < 0 {
		logrus.Errorf("pagination parameters are negative")
		return nil, errors.New("paginations parameters are negative(must be positiv)")
	}
	var query string
	
	logrus.Printf("parentID - %s",parentID)
	
	query = (`SELECT
			comment_id,
			post_id,
			parent_id,
			content,
			username,
			created_at,
			updated_at
		FROM
			comments c
		WHERE
			c.post_id = $1 AND parent_id = $2
		ORDER BY
			created_at DESC
		LIMIT $3 OFFSET $4;`)

		rows, err := r.db.Queryx(query, postID, parentID, *limit, *offset)
		if err != nil {
			logrus.Errorf("failed receive data from database:%s", err.Error())
			return nil, err
		}

		defer rows.Close()
		
		
		for rows.Next() {
			var comment model.Comment
			var updatedAt sql.NullString
			err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.Username, &comment.CreatedAt, &updatedAt)
			if err != nil {
				return nil, err
			}
	
			if updatedAt.Valid {
				comment.UpdatedAt = updatedAt.String
			}
	
			logrus.Printf("content - %s parent_id = %s", comment.Content, comment.ParentID)
			comments = append(comments, &comment)
		}
	

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return comments, nil
	
}

func (r *CommentPostgres) Update(ctx context.Context, input *model.UpdateComment) (string, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Content != nil {
		setValues = append(setValues, fmt.Sprintf("content=$%d", argId))
		args = append(args, *input.Content)
		argId++
	}

	setValues = append(setValues, "updated_at=NOW()")

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s  SET %s WHERE comment_id =$%d",
		commentsTable, setQuery, argId)
	args = append(args, input.ID)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		logrus.Errorf("failed to update comment in db: %s", err.Error())
		return "", err
	}

	return input.ID, nil
}

func (r *CommentPostgres) Delete(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE comment_id = $1", commentsTable)
	_, err := r.db.Exec(query, id)

	if err != nil {
		return false, errors.New("there are no comments with this id")
	}

	return true, nil
}
