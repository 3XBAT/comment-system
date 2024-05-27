package service

import (
	"context"

	"github.com/3XBAT/coments-system/graph/model"
	"github.com/3XBAT/coments-system/pkg/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(ctx context.Context, input *model.NewPost) (string, error) {
	return s.repo.Create(ctx, input)
}

func (s *PostService) Update(ctx context.Context, input *model.UpdatePost) (string, error) {
	return s.repo.Update(ctx, input)
}

func (s *PostService) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}

func (s *PostService) GetById(ctx context.Context, id string) (*model.Post, error){
	return s.repo.GetById(ctx, id)
}

func (s *PostService) GetAll(ctx context.Context) ([]*model.Post, error) {
	return s.repo.GetAll(ctx)
}