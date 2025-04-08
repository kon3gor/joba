package alert

import (
	"time"

	"github.com/kon3gor/joba/pkg/channel"
	"github.com/kon3gor/joba/pkg/formatter"
	"github.com/kon3gor/joba/pkg/scrapper"
	"github.com/kon3gor/joba/pkg/storage"
)

type AlertBuilder interface {
	ScrapUsing(scrapper.Scrapper) AlertBuilder
	Every(time.Duration) AlertBuilder
	FormatUsing(formatter.F) AlertBuilder
	SendInto(channel.C) AlertBuilder
	Build() *JobAlert

	SkipInitial(bool) AlertBuilder
}

type alertBuilder struct {
	alert *JobAlert
}

func NewAlert(id string, s storage.S) AlertBuilder {
	return &alertBuilder{
		alert: &JobAlert{
			ID:      id,
			storage: s,
		},
	}
}

func (ab *alertBuilder) ScrapUsing(s scrapper.Scrapper) AlertBuilder {
	ab.alert.scrapper = s
	return ab
}

func (ab *alertBuilder) Every(d time.Duration) AlertBuilder {
	ab.alert.checkPeriod = d
	return ab
}

func (ab *alertBuilder) FormatUsing(formatter formatter.F) AlertBuilder {
	ab.alert.foramtter = formatter
	return ab
}

func (ab *alertBuilder) SendInto(c channel.C) AlertBuilder {
	ab.alert.channel = c
	return ab
}

func (ab *alertBuilder) Build() *JobAlert {
	if ab.alert.scrapper == nil {
		panic("No scrapper provided for alert")
	}

	if ab.alert.checkPeriod == 0 {
		panic("No time period for alert provided")
	}

	if ab.alert.foramtter == nil {
		panic("No formatter for alert provided")
	}

	if ab.alert.channel == nil {
		panic("No channel for alert provided")
	}

	return ab.alert
}

func (ab *alertBuilder) SkipInitial(skip bool) AlertBuilder {
	ab.alert.skipInitial = skip
	return ab
}
