package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func IsSpamEmail(ctx context.Context, email string) (bool, error) {
	url := "https://" + os.Getenv("CHECK_SPAMMER_API") + email
	
	client := &http.Client{Timeout: time.Duration(time.Second * 5)}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body := struct {
		Spam bool `json:"spam"`
	}{}

	json.NewDecoder(res.Body).Decode(&body)

	return body.Spam, nil

}