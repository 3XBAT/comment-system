package service

import (
	"context"

	"github.com/3XBAT/coments-system/graph/model"
	"github.com/3XBAT/coments-system/pkg/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repos repository.Comment) *CommentService {
	return &CommentService{repo: repos}
}

func (s *CommentService) Create(ctx context.Context, input *model.NewComment) (string, error) {
	return s.repo.Create(ctx, input)
}

func (s *CommentService) GetById(ctx context.Context, id string) (*model.Comment, error) {
	return s.repo.GetById(ctx, id)
}

func (s *CommentService) GetPart(ctx context.Context, postID string, parentID string, limit *int, offset *int) ([]*model.Comment, error) {
	return s.repo.GetPart(ctx, postID, parentID, limit, offset ) 
}

func (s *CommentService) Update(ctx context.Context, input *model.UpdateComment) (string, error) {
	return s.repo.Update(ctx, input)
}

func (s *CommentService) Delete(ctx context.Context, id string) (bool, error){
	return s.repo.Delete(ctx, id)
}