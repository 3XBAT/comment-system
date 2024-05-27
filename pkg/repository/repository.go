package repository

import (
	"context"

	"github.com/3XBAT/coments-system/graph/model"
	"github.com/3XBAT/coments-system/pkg/repository/cache"
	"github.com/3XBAT/coments-system/pkg/repository/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)


type Post interface {
	Create(ctx context.Context, input *model.NewPost) (string, error)
	Update(ctx context.Context, input *model.UpdatePost) (string, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetById(ctx context.Context, id string) (*model.Post, error)
	GetAll(ctx context.Context) ([]*model.Post, error)
}

type Comment interface {
	Create(ctx context.Context, input *model.NewComment) (string, error)
	GetById(ctx context.Context, id string) (*model.Comment, error)
	GetPart(ctx context.Context, postID string, parentID string, limit *int, offset *int) ([]*model.Comment, error)
	Update(ctx context.Context, input *model.UpdateComment) (string, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type Repository struct {
	Post
	Comment
}

func NewRepositoryCache() *Repository {
	return &Repository{
		Post:    cache.NewPostCache(),
		Comment: cache.NewCommentCache(),
	}
}

func NewRepositoryPostgres(db *sqlx.DB) *Repository {
	logrus.Println("db is creating...")
	return &Repository{
		Post: postgres.NewPostPostgres(db),
		Comment: postgres.NewCommentPostgres(db),

	}
	
}
