package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

//go:generate mockgen -package=mock -source=cron.go -destination=$MOCK_FOLDER/cron.go Scheduler

const (
	timeFormat       = "15:04"
	defaultTime      = "00"
	defaultStartTime = "00:00"
)

const (
	minute = iota
	hour
	cmd
)

var timeLocation = time.UTC

var (
	errSettingCronJob = errors.New("failed setting cron job")
	errParsingTime    = errors.New("failed to parse time inputs")
)

// Cronjob services struct.
type Cronjob struct {
	minute string // minutes past the hour
	hour   string //  hour of the day
	cmd    string //  command to run
}

// Scheduler is the interface to interact with Cronjob.
type Scheduler interface {
	Schedule(offset string) (string, error)
}

// Start function starts cronjob and return time difference according to offset.
func (c *Cronjob) Schedule(offset string) (string, error) {
	timer, err := parseTimeFromCron(c.hour, c.minute)
	if err != nil {
		return "", err
	}

	err = c.createCronJob(timer)
	if err != nil {
		return "", err
	}

	parsedOffset, err := time.Parse(timeFormat, offset)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errSettingCronJob, err)
	}

	result := c.timeDifference(timer, parsedOffset)

	return result, nil
}

func (c *Cronjob) createCronJob(timer time.Time) error {
	s := gocron.NewScheduler(timeLocation).Days()

	switch {
	case runEveryTime(c.hour) && runEveryTime(c.minute):
		s = s.Every(1).Hour().Every(1).Minute().StartAt(timer)
	case runEveryTime(c.hour):
		s = s.Every(1).Hour().Every(timer.Minute()).Minute().StartAt(timer)
	case runEveryTime(c.minute):
		s = s.Every(timer.Hour).Hour().Every(1).Minute().StartAt(timer)
	default:
		s = s.Days().At(timer)
	}

	if _, err := s.Do(func() {}); err != nil {
		return fmt.Errorf("%w: %s", errSettingCronJob, err)
	}

	return nil
}

func (c *Cronjob) timeDifference(timer, offset time.Time) string {
	timeWithOffset := timer.Add(time.Hour*time.Duration(offset.Hour()) + time.Minute*time.Duration(offset.Minute()))

	minute := strconv.Itoa(timer.Minute())
	if minute == "0" {
		minute = defaultTime
	}

	switch {
	case runEveryTime(c.hour) && runEveryTime(c.minute):
		return fmt.Sprintf("%d:%d today - %s", timeWithOffset.Hour(), timeWithOffset.Minute(), c.cmd)
	case runEveryTime(c.hour):
		return fmt.Sprintf("%d:%d today - %s", timeWithOffset.Hour(), timeWithOffset.Minute(), c.cmd)
	case runEveryTime(c.minute):
		return fmt.Sprintf("%d:%s today - %s", timer.Hour(), minute, c.cmd)
	default:
		day := "today"
		if timeWithOffset.After(timer) {
			day = "tomorrow"
		}

		return fmt.Sprintf("%d:%s %s - %s", timer.Hour(), minute, day, c.cmd)
	}
}

func parseTimeFromCron(hour, minute string) (time.Time, error) {
	if hour == "*" {
		hour = defaultTime
	}

	if minute == "*" {
		minute = defaultTime
	}

	result, err := time.Parse(timeFormat, fmt.Sprintf("%s:%s", hour, minute))
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", errParsingTime, err)
	}

	return result, nil
}

// ParseSchedule function takes scheduler command line argument and parse it to Cronjob struct
func ParseSchedule(arg string) []Scheduler {
	cronjobs := make([]Cronjob, 0)

	for _, entry := range strings.Split(arg, "\n") {
		var c Cronjob

		for i, value := range strings.Fields(entry) {
			switch i {
			case minute:
				c.minute = value
			case hour:
				c.hour = value
			case cmd:
				c.cmd = value
			}
		}

		cronjobs = append(cronjobs, c)
	}

	schedulers := make([]Scheduler, len(cronjobs))
	for i := range cronjobs {
		schedulers[i] = Scheduler(&cronjobs[i])
	}

	return schedulers
}

func runEveryTime(input string) bool {
	return input == "*"
}
