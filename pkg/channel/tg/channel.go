package tg

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kon3gor/joba/pkg/channel"
	"github.com/kon3gor/joba/pkg/util"
	"golang.org/x/xerrors"
)

type Config struct {
	Token  string `yaml:"token"`
	ChatID string `yaml:"chat-id"`
}

type Channel struct {
	c Config
}

func NewChannel(c Config) channel.C {
	return &Channel{
		c: c,
	}
}

func (c *Channel) SendMessage(msg string) error {
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", util.GetRealValue(c.c.Token))

	data := url.Values{}
	data.Set("chat_id", util.GetRealValue(c.c.ChatID))
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
