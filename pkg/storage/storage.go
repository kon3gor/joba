package storage

import "context"

type S interface {
	Intersect(context.Context, string, []string) ([]string, error)
	Save(context.Context, string, []string) error
}
