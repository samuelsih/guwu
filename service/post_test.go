package service

import (
	"context"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/samuelsih/guwu/model"
	"github.com/stretchr/testify/assert"
)

func TestPostGetTimeline(t *testing.T) {
	ctx := context.Background()

	post := Post{
		DB: testDB,
	}

	guest := Guest{
		DB: testDB,
		SessionDB: testSessionDB,
	}

	loginUser := guest.Login(ctx, &GuestLoginIn{
		Email: "samuel@gmail.com",
		Password: "Akubohong123!",
	})

	t.Run(`Empty Timeline`, func(t *testing.T) {
		expected := PostTimelineOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		posts := post.Timeline(ctx)

		assert.Equal(t, expected.CommonResponse, posts.CommonResponse)
		assert.Nil(t, posts.Posts)
	})

	t.Run(`Showing Timeline`, func(t *testing.T) {
		_ = post.Insert(ctx, &PostInsertIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: loginUser.User.ID,
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},
	
			},
			Description: generateRandomString(350),
		})

		expected := PostTimelineOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := post.Timeline(ctx)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.Equal(t, 1, len(out.Posts))
	})
}

func TestPostInsert(t *testing.T) {
	ctx := context.Background()

	post := Post{
		DB: testDB,
	}

	guest := Guest{
		DB: testDB,
		SessionDB: testSessionDB,
	}

	loginUser := guest.Login(ctx, &GuestLoginIn{
		Email: "samuel@gmail.com",
		Password: "Akubohong123!",
	})

	t.Run(`Empty Struct`, func(t *testing.T) {
		in := PostInsertIn{}

		expected := PostInsertOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `description is required`,
			},
		}

		out := post.Insert(ctx, &in)

		assert.Equal(t, expected, out)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.Empty(t, out.Post)
	})

	t.Run(`Description Limit`, func(t *testing.T) {
		in := PostInsertIn{
			Description: generateRandomString(513),
		}

		expected := PostInsertOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `description is too long, max is 512 characters`,
			},
		}

		out := post.Insert(ctx, &in)

		assert.Equal(t, expected, out)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.Empty(t, out.Post)
	})

	t.Run(`Unknown User`, func(t *testing.T) {
		in := PostInsertIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: "123123123",
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},

			},
			Description: generateRandomString(100),
		}
		
		expected := PostInsertOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `unknown user`,
			},
		}

		out := post.Insert(ctx, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.Empty(t, out.Post)
	})

	t.Run(`Must Success`, func(t *testing.T) {
		in := PostInsertIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: loginUser.User.ID,
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},

			},
			Description: generateRandomString(100),
		}

		expected := PostInsertOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := post.Insert(ctx, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.NotEmpty(t, out.Post)
	})
}

func TestPostEdit(t *testing.T) {
	ctx := context.Background()

	post := Post{
		DB: testDB,
	}

	guest := Guest{
		DB: testDB,
		SessionDB: testSessionDB,
	}

	loginUser := guest.Login(ctx, &GuestLoginIn{
		Email: "samuel@gmail.com",
		Password: "Akubohong123!",
	})

	t.Run(`Unknown User`, func(t *testing.T) {
		insertedPost := post.Insert(ctx, &PostInsertIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: loginUser.User.ID,
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},
	
			},
			Description: generateRandomString(200),
		})

		expected := PostEditOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `unknown user`,
			},
		}

		in := PostEditIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: "09294836553",
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},
	
			},
			Description: "test 123",
		}

		out := post.Edit(ctx, insertedPost.Post.ID, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.Empty(t, out.Post)
	})

	t.Run(`Update Success`, func(t *testing.T) {
		insertedPost := post.Insert(ctx, &PostInsertIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: loginUser.User.ID,
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},
	
			},
			Description: generateRandomString(200),
		})

		expected := PostEditOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		in := PostEditIn{
			CommonRequest: CommonRequest{
				UserSession: model.Session{
					ID: loginUser.User.ID,
					Email: loginUser.User.Email,
					Username: loginUser.User.Username,
				},
	
			},
			Description: "test 123",
		}

		out := post.Edit(ctx, insertedPost.Post.ID, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.NotEmpty(t, out.Post)
	})
}



func generateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    
	b := make([]rune, n)
    
	for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    
	return string(b)
}

