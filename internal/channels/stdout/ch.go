package stdout

import (
	"fmt"
	"os"

	"github.com/kon3gor/joba/pkg/channel"
)

type stdOutChannel struct {
}

func NewChannel() channel.C {
	return &stdOutChannel{}
}

func (c *stdOutChannel) SendMessage(msg string) error {
	_, err := fmt.Fprintf(os.Stdout, "Got a message:\n%s", msg)
	return err
}
