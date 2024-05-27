package graph_test

import (
	"context"
	"errors"

	//"errors"
	"testing"

	"github.com/3XBAT/coments-system/graph/model"

	mock_service "github.com/3XBAT/coments-system/pkg/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResolver_CreatePost(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost, ctx context.Context, input *model.NewPost)

	testTable := []struct {
		name          string
		inputPost     *model.NewPost
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "OK_create_post",
			inputPost: &model.NewPost{
				Title:         "Test",
				Username:      "test",
				Content:       "test",
				AllowComments: true,
			},
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, input *model.NewPost) {
				s.EXPECT().Create(ctx, input).Return("generated-id", nil)
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			post := mock_service.NewMockPost(c)
			ctx := context.Background()
			testCase.mockBehavior(post, ctx, testCase.inputPost)

			responseID, err := post.Create(ctx, testCase.inputPost)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, responseID)
			}

		})
	}

}

func TestReslover_UpdatePost(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost, ctx context.Context, input *model.UpdatePost)

	testTable := []struct {
		name          string
		UpdatePost    *model.UpdatePost
		mockBehavior  mockBehavior
		expectedID    string
		expectedError error
	}{
		{
			name: "OK_update_post",
			UpdatePost: &model.UpdatePost{
				ID:            "92a84a90-140e-4d4c-a05e-e50ecc6216b8", //id существующего в бд поста
				Title:         stringPointer("test"),
				Content:       stringPointer("test"),
				Comment:       stringPointer("c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65"), //id существующего в бд комента
				AllowComments: boolPointer(true),
			},
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, input *model.UpdatePost) {
				s.EXPECT().Update(ctx, input).Return(input.ID, nil)
			},
			expectedID:    "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
			expectedError: nil,
		},
		{
			name: "Empty field",
			UpdatePost: &model.UpdatePost{
				ID:            "92a84a90-140e-4d4c-a05e-e50ecc6216b8", //id существующего в бд поста
				Content:       stringPointer("test"),
				Comment:       stringPointer("c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65"), //id существующего в бд комента
				AllowComments: boolPointer(true),
			},
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, input *model.UpdatePost) {
				s.EXPECT().Update(ctx, input).Return(input.ID, nil)
			},
			expectedID:    "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
			expectedError: nil,
		},
		{
			name: "New comment",
			UpdatePost: &model.UpdatePost{
				ID:      "92a84a90-140e-4d4c-a05e-e50ecc6216b8",                //id существующего в бд поста
				Comment: stringPointer("c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65"), //id существующего в бд комента
			},
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, input *model.UpdatePost) {
				s.EXPECT().Update(ctx, input).Return(input.ID, nil)
			},
			expectedID:    "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
			expectedError: nil,
		},
		{
			name: "Empty fields",
			UpdatePost: &model.UpdatePost{
				ID: "92a84a90-140e-4d4c-a05e-e50ecc6216b8", //id существующего в бд поста
			},
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, input *model.UpdatePost) {
				s.EXPECT().Update(ctx, input).Return("", errors.New("there is no data for changes"))
			},
			expectedID:    "",
			expectedError: errors.New("there is no data for changes"),
		},
		{
			name: "wrong id",
			UpdatePost: &model.UpdatePost{
				ID:            "92a84a9-4d4c-a05e-e50ecc6216b8", //id несуществующего в бд поста
				Title:         stringPointer("test"),
				Content:       stringPointer("test"),
				Comment:       stringPointer("c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65"), //id существующего в бд комента
				AllowComments: boolPointer(true),
			},
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, input *model.UpdatePost) {
				s.EXPECT().Update(ctx, input).Return("", errors.New("there is no post with this id"))
			},
			expectedID:    "",
			expectedError: errors.New("there is no post with this id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			post := mock_service.NewMockPost(c)
			ctx := context.Background()
			testCase.mockBehavior(post, ctx, testCase.UpdatePost)

			responseID, err := post.Update(ctx, testCase.UpdatePost)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError, err)
				assert.Equal(t, testCase.expectedID, responseID)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}
func boolPointer(b bool) *bool {
	return &b
}

func TestResolver_DeletePost(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost, ctx context.Context, id string)

	testTable := []struct {
		name           string
		id             string
		mockBehavior   mockBehavior
		expectedResult bool
		expectedError  error
	}{
		{
			name: "OK",
			id:   "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, id string) {
				s.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "wrong id",
			id:   "",
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, id string) {
				s.EXPECT().Delete(ctx, id).Return(false, errors.New("error while getting post by id"))
			},
			expectedResult: true,
			expectedError:  errors.New("error while getting post by id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			post := mock_service.NewMockPost(c)
			ctx := context.Background()
			testCase.mockBehavior(post, ctx, testCase.id)

			result, err := post.Delete(ctx, testCase.id)

			if err == nil {
				assert.Equal(t, testCase.expectedResult, result)
			} else {
				assert.Equal(t, testCase.expectedError, err)
			}

		})

	}
}

func TestResolver_GetPostById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost, ctx context.Context, id string)

	testTable := []struct {
		name          string
		id            string
		mockBehavior  mockBehavior
		expectedPost  *model.Post
		expectedError error
	}{
		{
			name: "OK",
			id:   "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, id string) {
				s.EXPECT().GetById(ctx, id).Return(&model.Post{ID: "92a84a90-140e-4d4c-a05e-e50ecc6216b8"}, nil)
			},
			expectedPost:  &model.Post{ID: "92a84a90-140e-4d4c-a05e-e50ecc6216b8"},
			expectedError: nil,
		},
		{
			name: "wrong id",
			id:   "",
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context, id string) {
				s.EXPECT().GetById(ctx, id).Return(nil, errors.New("there is no post with this id"))
			},
			expectedPost:  nil,
			expectedError: errors.New("there is no post with this id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			post := mock_service.NewMockPost(c)
			ctx := context.Background()
			testCase.mockBehavior(post, ctx, testCase.id)

			result, err := post.GetById(ctx, testCase.id)

			if err == nil {
				assert.Equal(t, testCase.expectedPost, result)
			} else {
				assert.Equal(t, testCase.expectedError, err)
			}
		})
	}
}

func TestResolver_GetAllPosts(t *testing.T) {
	type mockBehavior func(s *mock_service.MockPost, ctx context.Context)

	testTable := []struct {
		name          string
		id            string
		mockBehavior  mockBehavior
		expectedPosts []*model.Post
		expectedError error
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context) {
				s.EXPECT().GetAll(ctx).Return([]*model.Post{{ID: "92a84a90-140e-4d4c-a05e-e50ecc6216b8"}}, nil)
			},
			expectedPosts: []*model.Post{{ID: "92a84a90-140e-4d4c-a05e-e50ecc6216b8"}},
			expectedError: nil,
		},
		{
			name: "no post in data storage",
			mockBehavior: func(s *mock_service.MockPost, ctx context.Context) {
				s.EXPECT().GetAll(ctx).Return(nil, errors.New("there are no posts in the datastorage"))
			},
			expectedPosts: nil,
			expectedError: errors.New("there are no posts in the datastorage"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			post := mock_service.NewMockPost(c)
			ctx := context.Background()
			testCase.mockBehavior(post, ctx)

			posts, err := post.GetAll(ctx)

			if err == nil {
				assert.Equal(t, testCase.expectedPosts, posts)
			} else {
				assert.Equal(t, testCase.expectedError, err)
			}
		})
	}
}

func TestResolver_CreateComment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockComment, ctx context.Context, input *model.NewComment)

	testTable := []struct {
		name          string
		inputComment  *model.NewComment
		mockBehavior  mockBehavior
		exceptedId    string
		expectedError error
	}{
		{
			name: "OK with parent_id",
			inputComment: &model.NewComment{
				PostID:   "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
				ParentID: "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65",
				Content:  "test",
				Username: "test",
			},
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, input *model.NewComment) {
				s.EXPECT().Create(ctx, input).Return("genetated-id", nil)
			},
			expectedError: nil,
		},
		{
			name: "OK with parent_id",
			inputComment: &model.NewComment{
				PostID:   "92a84a90-140e-4d4c-a05e-e50ecc6216b8",
				ParentID: "",
				Content:  "test",
				Username: "test",
			},
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, input *model.NewComment) {
				s.EXPECT().Create(ctx, input).Return("genetated-id", nil)
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mock_service.NewMockComment(c)
			ctx := context.Background()
			testCase.mockBehavior(comment, ctx, testCase.inputComment)

			responseID, err := comment.Create(ctx, testCase.inputComment)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, responseID)
			}
		})
	}
}

func TestResolver_UpdateComment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockComment, ctx context.Context, input *model.UpdateComment)

	testTable := []struct {
		name          string
		UpdateComment *model.UpdateComment
		mockBehavior  mockBehavior
		expectedID    string
		expectedError error
	}{
		{
			name: "OK_update_comment",
			UpdateComment: &model.UpdateComment{
				ID:      "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65",
				Content: stringPointer("test"),
			},
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, input *model.UpdateComment) {
				s.EXPECT().Update(ctx, input).Return(input.ID, nil)
			},
			expectedID:    "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65",
			expectedError: nil,
		},
		{
			name: "empty content",
			UpdateComment: &model.UpdateComment{
				ID: "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65",
			},
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, input *model.UpdateComment) {
				s.EXPECT().Update(ctx, input).Return("", errors.New("there is no comment with this id"))
			},
			expectedID:    "",
			expectedError: errors.New("there is no comment with this id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mock_service.NewMockComment(c)
			ctx := context.Background()
			testCase.mockBehavior(comment, ctx, testCase.UpdateComment)

			responseID, err := comment.Update(ctx, testCase.UpdateComment)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError, err)
				assert.Equal(t, testCase.expectedID, responseID)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}

func TestResolver_DeleteComment(t *testing.T) {
	type mockBehavior func(s *mock_service.MockComment, ctx context.Context, id string)

	testTable := []struct {
		name           string
		id             string
		mockBehavior   mockBehavior
		expectedResult bool
		expectedError  error
	}{
		{
			name: "OK",
			id:   "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65",
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, id string) {
				s.EXPECT().Delete(ctx, id).Return(true, nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "wrong id",
			id:   "",
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, id string) {
				s.EXPECT().Delete(ctx, id).Return(false, errors.New("there are no comments with this id"))
			},
			expectedResult: false,
			expectedError:  errors.New("there are no comments with this id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mock_service.NewMockComment(c)
			ctx := context.Background()
			testCase.mockBehavior(comment, ctx, testCase.id)

			result, err := comment.Delete(ctx, testCase.id)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, result)
			}

		})

	}
}

func TestResolver_GetCommentById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockComment, ctx context.Context, id string)

	testTable := []struct {
		name            string
		id              string
		mockBehavior    mockBehavior
		expectedComment *model.Comment
		expectedError   error
	}{
		{
			name: "OK",
			id:   "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65",
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, id string) {
				s.EXPECT().GetById(ctx, id).Return(&model.Comment{ID: "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65"}, nil)
			},
			expectedComment: &model.Comment{ID: "c91ccfc9-e142-4c1c-8a1f-4c1d87fa9b65"},
			expectedError:   nil,
		},
		{
			name: "wrong id",
			id:   "",
			mockBehavior: func(s *mock_service.MockComment, ctx context.Context, id string) {
				s.EXPECT().GetById(ctx, id).Return(nil, errors.New("there are no comments with this id"))
			},
			expectedComment: nil,
			expectedError:   errors.New("there are no comments with this id"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mock_service.NewMockComment(c)
			ctx := context.Background()
			testCase.mockBehavior(comment, ctx, testCase.id)

			result, err := comment.GetById(ctx, testCase.id)

			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedComment, result)
			}
		})
	}
}

