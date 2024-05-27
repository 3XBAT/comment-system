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

type PostPostgres struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{db: db}
}

func (r *PostPostgres) Create(ctx context.Context, input *model.NewPost) (string, error) {

	NewPostId := (uuid.New().String())

	createPostQuery := fmt.Sprintf(`INSERT INTO %s (post_id, title, username, content, allow_comments, created_at)
	 VALUES ($1, $2, $3, $4, $5, $6)`, postsTable)

	created_at := time.Now().Format(time.RFC3339)

	_, err := r.db.Exec(createPostQuery, NewPostId, input.Title, input.Username, input.Content, input.AllowComments, created_at)

	if err != nil {
		logrus.Errorf("failed insert post into db: %s", err.Error())
		return "", err
	}

	return NewPostId, nil
}

func (r *PostPostgres) GetById(ctx context.Context, id string) (*model.Post, error) {

	query := (`SELECT
	p.post_id,
	p.title,
	p.username,
	p.content,
	p.allow_comments,
	p.created_at,
	p.updated_at,
	c.comment_id
FROM
	posts p
LEFT JOIN
	comments c ON p.post_id = c.post_id
WHERE
	p.post_id = $1`)

	rows, err := r.db.Queryx(query, id)
	if err != nil {
		logrus.Errorf("failed receive data from database:%s", err.Error())
		return nil, errors.New("there is no post with this id")
	}

	defer rows.Close()

	var post model.Post

	for rows.Next() {
		var commentID sql.NullString
		var updatedAt sql.NullString
		err := rows.Scan(&post.ID, &post.Title, &post.Username, &post.Content, &post.AllowComments, &post.CreatedAt,
			&updatedAt, &commentID)

		if err != nil {
			return nil, err
		}

		post.Comments = append(post.Comments, commentID.String)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(post.Comments) == 0 {
		logrus.Error("there are no objects with this id in the database")
		return nil, errors.New("there are no objects with this id in the database")
	}

	return &post, nil
}

func (r *PostPostgres) GetAll(ctx context.Context) ([]*model.Post, error) {

	query := `SELECT
		post_id,
		title,
		username,
		content,
		allow_comments,
		created_at,
		updated_at
	FROM
		posts
	ORDER BY
		created_at DESC;`

	rows, err := r.db.Queryx(query)
	if err != nil {
		logrus.Errorf("failed to retrieve data from database: %s", err.Error())
		return nil, errors.New("there are no posts in the datastorage")
	}

	var posts []*model.Post

	defer rows.Close()

	for rows.Next() {
		
		var UpdatedAt sql.NullString
		var post model.Post

		err := rows.Scan(
			&post.ID, &post.Title, &post.Username, &post.Content, &post.AllowComments, &post.CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		logrus.Errorf("error with scaning(GetAll posts)")
		return nil, err
	}

	return posts, nil
}

func (r *PostPostgres) Delete(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE post_id = $1", postsTable)
	_, err := r.db.Exec(query, id)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *PostPostgres) Update(ctx context.Context, input *model.UpdatePost) (string, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Content != nil {
		setValues = append(setValues, fmt.Sprintf("content=$%d", argId))
		args = append(args, *input.Content)
		argId++
	}

	if input.AllowComments != nil {
		setValues = append(setValues, fmt.Sprintf("allow_comments=$%d", argId))
		args = append(args, *input.AllowComments)
		argId++
	}

	setValues = append(setValues, "updated_at=NOW()")

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s p SET %s WHERE p.post_id =$%d",
		postsTable, setQuery, argId)
	args = append(args, input.ID)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		logrus.Errorf("failed to update post in db: %s", err.Error())
		return "", err
	}

	return input.ID, nil
}
