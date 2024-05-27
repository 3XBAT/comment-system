// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID        string   `json:"id"`
	PostID    string   `json:"postId"`
	ParentID  string   `json:"parentId"`
	Content   string   `json:"content"`
	Username  string   `json:"username"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Replies   []string `json:"replies,omitempty"`
}

type Mutation struct {
}

type NewComment struct {
	PostID   string `json:"postId"`
	ParentID string `json:"parentId"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

type NewPost struct {
	ID            *string `json:"id,omitempty"`
	Title         string  `json:"title"`
	Username      string  `json:"username"`
	Content       string  `json:"content"`
	AllowComments bool    `json:"allowComments"`
}

type Post struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Username      string   `json:"username"`
	Content       string   `json:"content"`
	Comments      []string `json:"comments,omitempty"`
	AllowComments bool     `json:"allowComments"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
}

type Query struct {
}

type UpdateComment struct {
	ID      string  `json:"id"`
	Content *string `json:"content,omitempty"`
}

type UpdatePost struct {
	ID            string  `json:"id"`
	Title         *string `json:"title,omitempty"`
	Content       *string `json:"content,omitempty"`
	Comment       *string `json:"comment,omitempty"`
	AllowComments *bool   `json:"allowComments,omitempty"`
}
