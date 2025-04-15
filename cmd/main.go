package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/kon3gor/gondor"
	"github.com/kon3gor/gondor/env"

	"github.com/kon3gor/joba/pkg"
	"github.com/kon3gor/joba/pkg/google"
	"github.com/kon3gor/joba/pkg/pg"
	"github.com/kon3gor/joba/pkg/tg"
)

var config struct {
	Alerts struct {
		Google struct {
			ID        string        `yaml:"id"`
			PageLimit int           `yaml:"page-limit"`
			Link      string        `yaml:"link"`
			Period    time.Duration `yaml:"period"`
		} `yaml:"google"`
	} `yaml:"alerts"`

	Postgres pg.Config `yaml:"db"`
	Telegram tg.Config `yaml:"telegram"`
}

func unmarshallDuration(d *time.Duration, b []byte) error {
	var err error
	*d, err = time.ParseDuration(string(b))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	gondor.RegisterStringHook(env.NewEnvHook())
	gondor.RegisterCustomUnmarshaler(unmarshallDuration)

	if err := gondor.Parse(&config, "config.yaml"); err != nil {
		log.Fatalln(err)
	}

	gc := google.NewScrapper(config.Alerts.Google.Link, config.Alerts.Google.PageLimit)
	st, closeF, err := pg.NewStorage(context.Background(), config.Postgres)
	if err != nil {
		log.Fatalln(err)
	}
	defer closeF()

	tgch := tg.NewChannel(config.Telegram)

	googleJobAlert := pkg.
		NewJobAlert(config.Alerts.Google.ID, st).
		ScrapUsing(gc).
		Every(config.Alerts.Google.Period).
		FormatUsing(SimpleFormatter{}).
		SendInto(tgch).
		SkipInitial().
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
