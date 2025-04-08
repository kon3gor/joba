package formatter

import (
	"github.com/kon3gor/joba/pkg/scrapper"
)

type F interface {
	Format([]scrapper.Result) string
}
