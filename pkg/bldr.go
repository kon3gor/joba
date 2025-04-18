package pkg

import (
	"time"
)

type AlertBuilder interface {
	ScrapUsing(Scrapper) AlertBuilder
	Every(time.Duration) AlertBuilder
	FormatUsing(Formatter) AlertBuilder
	SendInto(Channel) AlertBuilder
	Build() *JobAlert

	SkipInitial() AlertBuilder
}

type alertBuilder struct {
	alert *JobAlert
}

func NewJobAlert(id string, s Storage) AlertBuilder {
	return &alertBuilder{
		alert: &JobAlert{
			ID:      id,
			storage: s,
		},
	}
}

func (ab *alertBuilder) ScrapUsing(s Scrapper) AlertBuilder {
	ab.alert.scrapper = s
	return ab
}

func (ab *alertBuilder) Every(d time.Duration) AlertBuilder {
	ab.alert.checkPeriod = d
	return ab
}

func (ab *alertBuilder) FormatUsing(formatter Formatter) AlertBuilder {
	ab.alert.foramtter = formatter
	return ab
}

func (ab *alertBuilder) SendInto(c Channel) AlertBuilder {
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

func (ab *alertBuilder) SkipInitial() AlertBuilder {
	ab.alert.skipInitial = true
	return ab
}
