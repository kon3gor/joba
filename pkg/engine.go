package pkg

import (
	"context"
	"fmt"

	"golang.org/x/xerrors"
)

type Engine struct {
	alerts []*JobAlert
}

func NewEngine(alerts ...*JobAlert) *Engine {
	return &Engine{
		alerts: alerts,
	}
}

func (e *Engine) Run(ctx context.Context) error {
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

func (e *Engine) runIndividual(ctx context.Context, alert *JobAlert, errCh chan<- error) {
	err := alert.Run(ctx)
	if err != nil {
		errCh <- err
	}
}
