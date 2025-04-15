package tg

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"

	"github.com/kon3gor/joba/pkg"
)

type Config struct {
	Token  string `yaml:"token"`
	ChatID string `yaml:"chat-id"`
}

type Channel struct {
	c Config
}

func NewChannel(c Config) pkg.Channel {
	return &Channel{
		c: c,
	}
}

func (c *Channel) SendMessage(ctx context.Context, msg string) error {
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.c.Token)

	data := url.Values{}
	data.Set("chat_id", c.c.ChatID)
	data.Set("text", msg)

	// Send the HTTP POST request
	resp, err := http.PostForm(baseURL, data)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp)
		return xerrors.New("non 200 status code")
	}

	return nil
}
