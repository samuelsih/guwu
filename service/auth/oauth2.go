package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/service"
)

const (
	google = `google`
	github = `github`
)

type OAuthLoginOut struct {
	service.CommonResponse
	Link string `json:"string"`
}

func OAuthLogin(_ context.Context, provider string) OAuthLoginOut {
	var out OAuthLoginOut

	state, err := randomState(20)
	if err != nil {
		println(err)
		out.SetError(500, "can't generate oauth link state")
		return out
	}

	switch provider {
	case google:
		out.Link = config.GoogleConfig.AuthCodeURL(state)
		return out
	case github:
		out.Link = config.GithubConfig.AuthCodeURL(state)
		return out
	default:
		out.SetError(400, "unknown provider")
		return out
	}
}

type OAuthAuthorizeOut struct {
	service.CommonResponse
	User model.User `json:"user"`
}

func Authorize(ctx context.Context, provider, state, code string) OAuthAuthorizeOut {
	var out OAuthAuthorizeOut
	var user model.User

	switch provider {
	case google:
		token, err := config.GoogleConfig.Exchange(ctx, code)
		if err != nil {
			out.SetError(400, err.Error())
			return out
		}

		if !token.Valid() {
			out.SetError(500, `invalid token received from provider`)
			return out
		}

		client := config.GoogleConfig.Client(ctx, token)
		resp, err := client.Get("https://api.github.com/user")

		if err != nil {
			out.SetError(500, err.Error())
			return out
		}
		
		user, err = extract(provider, resp)
		if err != nil {
			out.SetError(500, err.Error())
			return out
		}

	case github:
		token, err := config.GithubConfig.Exchange(ctx, code)
		if err != nil {
			out.SetError(400, err.Error())
			return out
		}

		req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
		if err != nil {
			out.SetError(500, `cant make http request`)
			return out
		}

		req.Header.Add("Authorization", "Bearer " + token.AccessToken)

		client := config.GithubConfig.Client(ctx, token)
		resp, err := client.Do(req)

		if err != nil {
			out.SetError(500, err.Error())
			return out
		}
		
		user, err = extract(provider, resp)
		if err != nil {
			out.SetError(500, err.Error())
			return out
		}

		if user.Email == "" {
			user.Email, err = getGithubPrivateEmail(client, token.AccessToken)
			if err != nil {
				out.SetError(406, err.Error())
				return out
			}
		}

	default:
		out.SetError(400, "unknown provider")
		return out
	}

	if user.CreatedAt == (time.Time{}) {
		user.CreatedAt = time.Now()
	}

	out.User = user
	out.SetOK()
	return out
}

func randomState(n int) (string, error) {
    data := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, data); err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(data), nil
}

func extract(provider string, r *http.Response) (model.User, error) {
	if r == nil {
		return model.User{}, fmt.Errorf("Response is nil on provider %s", provider)
	}
	
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return model.User{}, fmt.Errorf("Provider %s is responded %d while trying to get user information", provider, r.StatusCode)
	}

	var user model.User

	switch provider {
	case google:
		responseBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return model.User{}, err
		}
	
		if err := json.Unmarshal(responseBytes, &user); err != nil {
			return user, err
		}
	
		return user, nil

	case github:
		githubStruct := struct{
			ID int    `json:"id"`
			Name string  `json:"name"`
			Email string `json:"email"`
		}{}

		responseBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return model.User{}, err
		}
	
		if err := json.Unmarshal(responseBytes, &githubStruct); err != nil {
			return user, err
		}

		user.ID = strconv.Itoa(githubStruct.ID)
		user.Name = githubStruct.Name
		user.Email = githubStruct.Email
	}

	return user, nil
}

func getGithubPrivateEmail(client *http.Client, token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer " + token)

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var mailList []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	err = json.NewDecoder(res.Body).Decode(&mailList)
	
	if err != nil {
		return "", err
	}

	for _, v := range mailList {
		if v.Primary && v.Verified {
			return v.Email, nil
		}
	}

	return "", fmt.Errorf("Email must be verified at github")
}