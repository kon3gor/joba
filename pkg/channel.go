package pkg

import "context"

type Channel interface {
	SendMessage(context.Context, string) error
}
