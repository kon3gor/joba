package pkg

type Formatter interface {
	Format([]ScrapResult) string
}
