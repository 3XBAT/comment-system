package cache

import (
	"context"
	"errors"
	"time"

	"github.com/3XBAT/coments-system/graph/model"
	
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostCache struct {
	Cache map[string]model.Post
}

func NewPostCache() *PostCache  {
	return &PostCache{Cache: map[string]model.Post{}}
}

func (c *PostCache) Create(ctx context.Context, input *model.NewPost) (string, error) {
	
	newInput := model.Post{
		ID:            uuid.New().String(),
		Title:         input.Title,
		Username:      input.Username,
		Content:       input.Content,
		Comments:      nil,
		AllowComments: input.AllowComments,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	c.Cache[newInput.ID] = newInput
	return newInput.ID, nil
}

func (c *PostCache) Update(ctx context.Context, input *model.UpdatePost) (string, error) {
	if !constainsKeyPost(input.ID, c.Cache){
		logrus.Errorf("there is no post with this id")
		return "", errors.New("there is no post with this id")
	}
	
	post := c.Cache[input.ID];

	if input.AllowComments != nil {
		post.AllowComments = *input.AllowComments
	}
	
	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	if input.Comment != nil {
		post.Comments = append(post.Comments, *input.Comment)
	}

	post.UpdatedAt = time.Now().Format(time.RFC3339)
	

	c.Cache[input.ID] = post
	
	return input.ID, nil

}

func (c *PostCache) Delete(ctx context.Context, id string) (bool, error) {
	
	if !constainsKeyPost(id, c.Cache){
		logrus.Fatalf("There is no post with this id")
		return false, errors.New("there is no post with this id")
	}
	
	delete(c.Cache,id)
	
	return true, nil
}

func (c *PostCache) GetById(ctx context.Context, id string) (*model.Post, error) {
	var item model.Post
	
	if !constainsKeyPost(id, c.Cache){
		logrus.Errorf("there is no post with this id")
		return &item, errors.New("there is no post with this id")
	}

	item = c.Cache[id]

	return &item, nil
}

func (c *PostCache) GetAll(ctx context.Context) ([]*model.Post, error) {
	if len(c.Cache) == 0 {
		logrus.Errorf("there are no posts in the data storage(cache)")
		return nil, errors.New("there are no posts in the data storage")
	}

	var posts []*model.Post

	for _, post := range c.Cache{
		posts = append(posts, &post)
	}

	return posts, nil
}
