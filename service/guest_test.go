package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuestRegister(t *testing.T) {
	ctx := context.Background()

	guest := Guest{
		DB: testDB,
		SessionDB: testSessionDB,
	}

	t.Run(`Register Empty struct`, func(t *testing.T) {
		in := GuestRegisterIn{}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: errEmailRequired.Error(),
			},
		}

		out := guest.Register(ctx, &in)

		assert.Equal(t, expected, out)
		assert.Empty(t, out.User)
	})

	t.Run(`Insert Success`, func(t *testing.T) {
		in := GuestRegisterIn{
			Email: "samuel02@gmail.com",
			Username: "SamuelSamuelSamuel",
			Password: "Akubohong123!",
		}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := guest.Register(ctx, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.NotEmpty(t, out.User)
		assert.NotEmpty(t, out.SessionID)
	})

	t.Run(`Duplicate Insert`, func(t *testing.T) {
		in := GuestRegisterIn{
			Email: "samuel90@gmail.com",
			Username: "HayangUlinAing",
			Password: "GajadiBohong123!",
		}

		in2 := GuestRegisterIn{
			Email: "samuel90@gmail.com",
			Username: "HayangUlinAing",
			Password: "GajadiBohong123!",
		}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `email already exists`,
			},
		}

		_ = guest.Register(ctx, &in)
		out := guest.Register(ctx, &in2)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	guest := Guest{
		DB: testDB,
		SessionDB: testSessionDB,
	}

	t.Run(`Login Empty struct`, func(t *testing.T) {
		in := GuestLoginIn{}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `invalid email`,
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.Empty(t, out.User)
	})

	t.Run(`Login Empty Password`, func(t *testing.T){
		in := GuestLoginIn{
			Email: "aiaiai@gmail.com",
		}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `invalid password`,
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.Empty(t, out.User)
	})

	t.Run(`Login Email Not Found`, func(t *testing.T) {
		in := GuestLoginIn{
			Email: "bukansamuel@gmail.com",
			Password: "Passwordajah",
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, http.StatusBadRequest, out.StatusCode)
		assert.Empty(t, out.User)
	})

	t.Run(`Login With Wrong Password`, func(t *testing.T) {
		_ = guest.Register(ctx, &GuestRegisterIn{
			Email: "samuel911@gmail.com",
			Username: "HayangUlinAing",
			Password: "GajadiBohong123!",
		})

		in := GuestLoginIn{
			Email: "samuel911@gmail.com",
			Password: "GajadiBohong123",
		}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `email or password does not match`,
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
		assert.Empty(t, out.User)
	})

	t.Run(`Success`, func(t *testing.T){
		in := GuestLoginIn{
			Email: "samuel911@gmail.com",
			Password: "GajadiBohong123!",
		}

		expected := GuestLoginOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := guest.Login(ctx, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.NotEmpty(t, out.User)
		assert.NotEmpty(t, out.User.ID)
		assert.NotEmpty(t, out.SessionID)
		assert.Equal(t, out.User.Email, "samuel911@gmail.com")
		assert.Equal(t, out.User.Username, "HayangUlinAing")
	})
}

func TestLogout(t *testing.T) {
	ctx := context.Background()

	var sessionID string

	guest := Guest{
		DB: testDB,
		SessionDB: testSessionDB,
	}

	t.Run(`Insert Success`, func(t *testing.T) {
		in := GuestRegisterIn{
			Email: "testuser@gmail.com",
			Username: "TestUser",
			Password: "Akubohong123!",
		}

		expected := GuestRegisterOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := guest.Register(ctx, &in)

		assert.Equal(t, expected.StatusCode, out.StatusCode)
		assert.Equal(t, expected.Msg, out.Msg)
		assert.NotEmpty(t, out.User)
		assert.NotEmpty(t, out.SessionID)

		sessionID = out.SessionID
	})

	t.Run(`Empty Session`, func(t *testing.T){
		expected := GuestLogoutOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg: `user not found`,
			},
		}

		out := guest.Logout(ctx, "")

		assert.Equal(t, expected.CommonResponse, out.CommonResponse)
	})

	t.Run(`Logout Must Success`, func(t *testing.T){
		expected := GuestLogoutOut{
			CommonResponse: CommonResponse{
				StatusCode: http.StatusOK,
				Msg: "OK",
			},
		}

		out := guest.Logout(ctx, sessionID)

		assert.Equal(t, out.CommonResponse, expected.CommonResponse)
	})
}