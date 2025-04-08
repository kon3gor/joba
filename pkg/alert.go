package pkg

import (
	"context"
	"slices"
	"time"
)

type JobAlert struct {
	ID          string
	scrapper    Scrapper
	channel     Channel
	checkPeriod time.Duration
	storage     Storage
	foramtter   Formatter

	skipInitial bool
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

	filteredResults := make([]ScrapResult, 0, len(filteredIds))
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

	if err := ja.channel.SendMessage(ctx, ja.foramtter.Format(filteredResults)); err != nil {
		return err
	}

save:
	return ja.storage.Save(ctx, ja.ID, filteredIds)
}
