package service

import (
	"context"

	"github.com/3XBAT/coments-system/graph/model"
	"github.com/3XBAT/coments-system/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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

type Service struct {
	Post
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Post: NewPostService(repos.Post),
		Comment: NewCommentService(repos.Comment),
	}
}
