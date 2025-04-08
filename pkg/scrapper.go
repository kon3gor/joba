package pkg

import "fmt"

type ScrapResult interface {
	fmt.Stringer

	GetID() string
	GetLink() string
}

type Scrapper interface {
	Scrap() ([]ScrapResult, error)
}
