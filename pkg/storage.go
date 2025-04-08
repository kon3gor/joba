package pkg

import "context"

type Storage interface {
	Intersect(context.Context, string, []string) ([]string, error)
	Save(context.Context, string, []string) error
}
