package process

import (
	"fmt"
	"time"

	"github.com/warmans/dbr"
	"go.uber.org/zap"
)

type Runner struct {
	processor Processor
	interval  time.Duration
	logger    *zap.Logger
}

func (r *Runner) Run(db *dbr.Session) {

	logger := r.logger.With(zap.String("processor", fmt.Sprintf("%T", r.processor)))

	logger.Info(fmt.Sprintf("Starting processor. Updating every %d minutes", int64(r.interval / time.Minute)))
	for {
		startTime := time.Now()
		if err := r.processor.Update(db); err != nil {
			logger.Error(fmt.Sprintf("Processor failed to complete update with error: %s", err.Error()))
		}

		runDuration := time.Since(startTime)
		logger.Info(fmt.Sprintf("Ran in %d seconds", int64(runDuration / time.Second)))

		if waitTime := r.interval - runDuration; waitTime > 0 {
			time.Sleep(waitTime)
		}
	}
}

type Processor interface {
	Update(db *dbr.Session) error
}

func GetActivityRunner(interval time.Duration, logger *zap.Logger) *Runner {
	return &Runner{processor: &Activity{}, interval: interval, logger: logger}
}

type Activity struct{}

func (p *Activity) Update(db *dbr.Session) error {

	if _, err := db.Exec(
		`UPDATE performer SET activity =  (
			SELECT SUM(1) as num
			FROM event_performer  ep
			LEFT JOIN event e ON ep.event_id = e.id
			WHERE ep.performer_id = performer.id
			AND e.date > date('now', '-1 month')
		)`,
	); err != nil {
		return err
	}

	_, err := db.Exec(
		`UPDATE venue SET activity =  (
			SELECT SUM(1) as num
			FROM event e
			WHERE e.venue_id = venue.id
			AND e.date > date('now', '-1 month')
		)`,
	)

	return err
}
