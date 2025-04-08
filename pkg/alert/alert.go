package alert

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/kon3gor/joba/pkg/channel"
	"github.com/kon3gor/joba/pkg/formatter"
	"github.com/kon3gor/joba/pkg/scrapper"
	"github.com/kon3gor/joba/pkg/storage"
)

type JobAlert struct {
	ID          string
	scrapper    scrapper.Scrapper
	channel     channel.C
	checkPeriod time.Duration
	storage     storage.S
	foramtter   formatter.F

	skipInitial bool
}

func NewJobAlert(
	scrap scrapper.Scrapper,
	every time.Duration,
	into channel.C,
	storage storage.S,
	foramtter formatter.F,
) (*JobAlert, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}

	return &JobAlert{
		ID:          id.String(),
		scrapper:    scrap,
		channel:     into,
		checkPeriod: every,
		storage:     storage,
		foramtter:   foramtter,
	}, nil
}

func (ja *JobAlert) Run(ctx context.Context) error {
	ticker := time.NewTicker(ja.checkPeriod)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ticker.Stop()
			err := ja.perform(ctx)
			if err != nil {
				ticker.Stop()
				return err
			}
			ticker.Reset(ja.checkPeriod)
		}
	}
}

func (ja *JobAlert) perform(ctx context.Context) error {
	results, err := ja.scrapper.Scrap()
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(results))
	for _, res := range results {
		ids = append(ids, res.GetID())
	}

	filteredIds, err := ja.storage.Intersect(ctx, ja.ID, ids)
	if err != nil {
		return err
	}

	filteredResults := make([]scrapper.Result, 0, len(filteredIds))
	for _, res := range results {
		if slices.Contains(filteredIds, res.GetID()) {
			filteredResults = append(filteredResults, res)
		}
	}

	if len(filteredResults) == 0 {
		return nil
	}

	if ja.skipInitial {
		ja.skipInitial = false
		goto save
	}

	if err := ja.channel.SendMessage(ja.foramtter.Format(filteredResults)); err != nil {
		return err
	}

save:
	return ja.storage.Save(ctx, ja.ID, filteredIds)
}
