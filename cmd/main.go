package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/kon3gor/gondor"

	"github.com/kon3gor/joba/pkg"
	"github.com/kon3gor/joba/pkg/google"
	"github.com/kon3gor/joba/pkg/pg"
	"github.com/kon3gor/joba/pkg/tg"
)

var config struct {
	Google struct {
		ID        string `yaml:"id"`
		PageLimit int    `yaml:"page-limit"`
		Link      string `yaml:"link"`
	} `yaml:"google"`

	Postgres pg.Config `yaml:"db"`
	Telegram tg.Config `yaml:"telegram"`
}

func main() {
	if err := gondor.Parse(&config, "config.yaml"); err != nil {
		log.Fatalln(err)
	}

	gc := google.NewScrapper(config.Google.Link, config.Google.PageLimit)
	st, closeF, err := pg.NewStorage(context.Background(), config.Postgres)
	if err != nil {
		log.Fatalln(err)
	}
	defer closeF()

	tgch := tg.NewChannel(config.Telegram)

	googleJobAlert := pkg.
		NewJobAlert(config.Google.ID, st).
		ScrapUsing(gc).
		Every(5 * time.Second).
		FormatUsing(SimpleFormatter{}).
		SendInto(tgch).
		SkipInitial(false).
		Build()

	engine := pkg.NewEngine(googleJobAlert)
	if err := engine.Run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

type SimpleFormatter struct {
}

func (sf SimpleFormatter) Format(r []pkg.ScrapResult) string {
	var msg strings.Builder

	msg.WriteString("New jobs!!!!\n")
	for _, res := range r {
		msg.WriteString(res.String())
		msg.WriteByte('\n')
	}

	return msg.String()
}
