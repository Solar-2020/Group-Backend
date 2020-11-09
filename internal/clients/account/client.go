package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Solar-2020/Account-Backend/pkg/models"
	"net/http"
	"strconv"
)

type Client interface {
	GetUserIDByEmail(email string) (user int, err error)
	GetUserByID(userID int) (user models.User, err error)
}

type client struct {
	host   string
	secret string
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret}
}

type httpError struct {
	Error string `json:"error"`
}

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Surname   string `json:"surname"`
	AvatarURL string `json:"avatarURL"`
}

func (c *client) GetUserIDByEmail(email string) (userID int, err error) {
	req, err := http.NewRequest(http.MethodGet, c.host+fmt.Sprintf("/api/internal/account/by-email/%s", email), nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response User
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response.ID, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return userID, errors.New(httpErr.Error)
	default:
		return userID, errors.New("Unexpected Server Error")
	}
}

func (c *client) GetUserByID(userID int) (user models.User, err error) {
	req, err := http.NewRequest(http.MethodGet, c.host+fmt.Sprintf("/api/internal/account/by-user/%s", strconv.Itoa(userID)), nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response models.User
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return user, errors.New(httpErr.Error)
	default:
		return user, errors.New("Unexpected Server Error")
	}
}
