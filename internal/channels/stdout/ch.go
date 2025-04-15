package stdout

import (
	"context"
	"fmt"
	"os"

	"github.com/kon3gor/joba/pkg"
)

type stdOutChannel struct {
}

func NewChannel() pkg.Channel {
	return &stdOutChannel{}
}

func (c *stdOutChannel) SendMessage(ctx context.Context, msg string) error {
	_, err := fmt.Fprintf(os.Stdout, "Got a message:\n%s", msg)
	return err
}
