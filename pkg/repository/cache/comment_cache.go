package cache

import (
	"context"
	"errors"
	"time"

	"github.com/3XBAT/coments-system/graph/model"
	

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CommentCache struct {
	Cache map[string]model.Comment
}

func NewCommentCache() *CommentCache {
	return &CommentCache{Cache: map[string]model.Comment{}}
	
}

func (c *CommentCache) Create(ctx context.Context, input *model.NewComment) (string, error) {

	if len(input.Content) > 2000 {
		logrus.Errorf("the length of comment over then 2000 symbols")
		return "", errors.New("the length of comment over then 2000 symbols")
	}


	newComment := model.Comment{
		ID:        uuid.New().String(),
		PostID:    input.PostID,
		ParentID:  input.ParentID,
		Content:   input.Content,
		Username:  input.Username,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	c.Cache[newComment.ID] = newComment
	return newComment.ID, nil
}

func (c *CommentCache) GetById(ctx context.Context, id string) (*model.Comment, error) {
	var comment model.Comment

	if !constainsKeyComment(id, c.Cache) {
		logrus.Errorf("There is are comments with this id")
		return nil, errors.New("there are no comments with this id")
	}

	comment = c.Cache[id]

	return &comment, nil
}

func (c *CommentCache) GetPart(ctx context.Context, postID string, parentID string, limit *int, offset *int) ([]*model.Comment, error) {
	var comments []*model.Comment

	if *limit <= 0 || *offset < 0 {
		logrus.Errorf("pagination parameters are negative")
		return nil, errors.New("paginations parameters are negative(must be positiv)")
	}

	counter := 0
	if parentID == ""{
		for _, comment := range c.Cache {
			if counter >= *offset {
				if counter < *limit + *offset - 1 {
					if comment.PostID == postID {
						comments = append(comments, &comment)
					}
				}
			}

			counter++
		}
	}else{
		for _, comment := range c.Cache{
			if counter >= *offset{
				if counter < *limit{
					if comment.ParentID == parentID{
						comments = append(comments, &comment)
					}
				}
			}
		}
	}

	if len(comments) == 0 {
		logrus.Errorf("there are no comments with this id in the data storage")
		return nil, errors.New("there are no comments with this id in the data storage")
	}

	return comments, nil
}

func (c *CommentCache) Update(ctx context.Context, input *model.UpdateComment) (string, error) {
	if !constainsKeyComment(input.ID, c.Cache) {
		logrus.Fatalf("there is no comment with this id")
		return "", errors.New("there is no comment with this id")
	}

	comment := c.Cache[input.ID]

	if input.Content != nil {
		comment.Content = *input.Content
	}

	comment.UpdatedAt = time.Now().Format(time.RFC3339)

	c.Cache[input.ID] = comment

	return input.ID, nil
}

func (c *CommentCache) Delete(ctx context.Context, id string) (bool, error) {
	if !constainsKeyComment(id, c.Cache) {
		logrus.Fatalf("there is no comment with this id")
		return false, errors.New("there is no comment with this id")
	}

	delete(c.Cache, id)

	return true, nil
}
