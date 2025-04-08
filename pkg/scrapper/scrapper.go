package scrapper

import "fmt"

type Result interface {
	fmt.Stringer

	GetID() string
	GetLink() string
}

type Scrapper interface {
	Scrap() ([]Result, error)
}
