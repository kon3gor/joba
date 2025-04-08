package engine

import (
	"context"
	"fmt"

	"github.com/kon3gor/joba/pkg/alert"
	"golang.org/x/xerrors"
)

type E struct {
	alerts []*alert.JobAlert
}

func NewEngine(alerts ...*alert.JobAlert) *E {
	return &E{
		alerts: alerts,
	}
}

func (e *E) Run(ctx context.Context) error {
	alertErrs := make(chan error, len(e.alerts))
	defer close(alertErrs)

	activeAlerts := len(e.alerts)
	for _, alert := range e.alerts {
		go func() {
			e.runIndividual(ctx, alert, alertErrs)
		}()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-alertErrs:
			// TODO(kon3gor): Better logging for errors
			activeAlerts--

			fmt.Println(err)
		}

		if activeAlerts == 0 {
			break
		}
	}

	return xerrors.New("all alerts failed")
}

func (e *E) runIndividual(ctx context.Context, alert *alert.JobAlert, errCh chan<- error) {
	err := alert.Run(ctx)
	if err != nil {
		errCh <- err
	}
}
