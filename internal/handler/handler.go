package handler

import (
	"errors"
	"flag"
	"fmt"

	"github.com/TestardR/cron_scheduler/internal/cron"
	"github.com/TestardR/cron_scheduler/internal/logger"
)

const (
	scheduleUsage = `Please enter a schedule in the following format:
30 1 /bin/run_me_daily
45 * /bin/run_me_hourly
* * /bin/run_me_every_minute
* 19 /bin/run_me_sixty_times
`
	offsetUsage = "Please enter a offset value in the following format: hh:mm"
)

var (
	errMissingArgument = errors.New("missing command line argument")
	errFailedCronjob   = errors.New("an error occurred on a cronjob")
)

// Handler services struct.
type Handler struct {
	log      logger.Logger
	cronjobs []cron.Scheduler
	offset   string
}

// Start function starts Handler service.
func (h *Handler) Start() error {
	h.log.Info("scheduling cronjobs started")

	for _, cronjob := range h.cronjobs {
		output, err := cronjob.Schedule(h.offset)
		if err != nil {
			err := fmt.Errorf("%w:%s", errFailedCronjob, err)
			h.log.Error(err)

			return err
		}

		fmt.Println(output)
	}

	h.log.Info("scheduling cronjobs finished")

	return nil
}

// New function iinitializes Handler service.
func New(log logger.Logger) (*Handler, error) {
	s := flag.String("schedule", "", scheduleUsage)
	t := flag.String("offset", "", offsetUsage)
	flag.Parse()

	if *s == "" || *t == "" {
		return nil, errMissingArgument
	}

	sch := Handler{
		log:      log,
		cronjobs: cron.ParseSchedule(*s),
		offset:   *t,
	}

	sch.log.Info("initialized application successfully")

	return &sch, nil
}
