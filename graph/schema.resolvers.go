package graph

import (
	"context"
	"errors"

	"github.com/3XBAT/coments-system/graph/model"
	"github.com/sirupsen/logrus"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input *model.NewPost) (string, error) {
	return r.Resolver.service.Post.Create(ctx, input)
}

func (r *mutationResolver) UpdatePost(ctx context.Context, input *model.UpdatePost) (string, error) {
	if input.AllowComments == nil && input.Content == nil && input.Title == nil {
		return "", errors.New("there is no data for changes")
	}

	if !*input.AllowComments { // если запрещаем комментарии, то оставленные раннее удаляются
		post, err := r.Resolver.service.Post.GetById(ctx, input.ID)
		if err != nil {
			return "", err
		}

		for _, comment_id := range post.Comments {
			r.Resolver.service.Comment.Delete(ctx, comment_id)
		}
	}

	return r.Resolver.service.Post.Update(ctx, input)
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	post, err := r.Resolver.service.Post.GetById(ctx, id)
	if err != nil {
		logrus.Errorf("error while getting post by id:%s", err.Error())
		return false, errors.New("error while getting post by id")
	}

	maxLenOfComments := len(post.Comments)
	offsetForDelete := 0

	comments, err := r.Resolver.service.Comment.GetPart(ctx, id, id, &maxLenOfComments, &offsetForDelete)
	if err != nil {
		logrus.Errorf("error while getting comments: %s", err.Error())
		return r.Resolver.service.Post.Delete(ctx, id)
	}

	for _, comment := range comments {
		r.Resolver.service.Comment.Delete(ctx, comment.ID)
	}

	return r.Resolver.service.Post.Delete(ctx, id)
}

func (r *mutationResolver) CreateComment(ctx context.Context, input *model.NewComment) (string, error) {
	_, err := r.Resolver.service.Post.GetById(ctx, input.PostID)
	if err != nil {
		logrus.Errorf("there is no post with this id")
		return "", errors.New("there is no post with this id")
	}

	commentId, err := r.Resolver.service.Comment.Create(ctx, input)
	if err != nil {
		logrus.Errorf("Erorr while creating comment: %s", err.Error())
		return "", err
	}

	comment, err := r.Resolver.service.Comment.GetById(ctx, commentId)
	if err != nil {
		logrus.Errorf("error while receiving a comment %s", err.Error())
		return "", err
	}

	updPost := model.UpdatePost{
		ID:      comment.PostID,
		Comment: &commentId,
	}

	r.Resolver.service.Post.Update(ctx, &updPost)

	return commentId, nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, input *model.UpdateComment) (string, error) {
	return r.Resolver.service.Comment.Update(ctx, input)
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	return r.Resolver.service.Comment.Delete(ctx, id)
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.Resolver.service.Post.GetAll(ctx)
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return r.Resolver.service.Post.GetById(ctx, id)
}

func (r *queryResolver) Comments(ctx context.Context, postID string, parentID string, limit *int, offset *int) ([]*model.Comment, error) {
	return r.Resolver.service.Comment.GetPart(ctx, postID, parentID, limit, offset)
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id string) (*model.Comment, error) {
	return r.Resolver.service.Comment.GetById(ctx, id)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
