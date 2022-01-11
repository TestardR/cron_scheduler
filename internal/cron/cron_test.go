package cron

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cronJobTestCase struct {
	c      Cronjob
	offset string
	out    string
	err    error
}

func TestStart(t *testing.T) {
	t.Parallel()
	tests := map[string]cronJobTestCase{
		"defined-minute-and-hour-offset-after":  CaseDefinedMinuteAndHourOffsetAfter(t),
		"defined-minute-and-hour-offset-before": CaseDefinedMinuteAndHourOffsetBefore(t),
		"undefined-minute-and-hour":             CaseUndefinedMinutAndHour(t),
		"defined-minute-and-undefined-hour":     CaseDefinedMinutAndUnDefinedHour(t),
		"undefined-minute-and-defined-hour":     CaseUndefinedMinuteAndDefinedHour(t),
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			result, err := tc.c.Schedule(tc.offset)
			require.NoError(t, err)
			assert.Equal(t, result, tc.out)
		})
	}
}

func CaseDefinedMinuteAndHourOffsetAfter(t *testing.T) cronJobTestCase {
	c := Cronjob{
		minute: "30",
		hour:   "1",
		cmd:    "/bin/run_me_daily",
	}

	return cronJobTestCase{
		offset: "16:10",
		c:      c,
		out:    "1:30 tomorrow - /bin/run_me_daily",
		err:    nil,
	}
}

func CaseDefinedMinuteAndHourOffsetBefore(t *testing.T) cronJobTestCase {
	c := Cronjob{
		minute: "30",
		hour:   "1",
		cmd:    "/bin/run_me_daily",
	}

	return cronJobTestCase{
		offset: "00:00",
		c:      c,
		out:    "1:30 today - /bin/run_me_daily",
		err:    nil,
	}
}

func CaseUndefinedMinutAndHour(t *testing.T) cronJobTestCase {
	c := Cronjob{
		minute: "*",
		hour:   "*",
		cmd:    "/bin/run_me_every_minute",
	}

	return cronJobTestCase{
		offset: "16:10",
		c:      c,
		out:    "16:10 today - /bin/run_me_every_minute",
		err:    nil,
	}
}

func CaseDefinedMinutAndUnDefinedHour(t *testing.T) cronJobTestCase {
	c := Cronjob{
		minute: "45",
		hour:   "*",
		cmd:    "/bin/run_me_hourly",
	}

	return cronJobTestCase{
		offset: "16:10",
		c:      c,
		out:    "16:55 today - /bin/run_me_hourly",
		err:    nil,
	}
}

func CaseUndefinedMinuteAndDefinedHour(t *testing.T) cronJobTestCase {
	c := Cronjob{
		minute: "*",
		hour:   "19",
		cmd:    "/bin/run_me_sixty_times",
	}

	return cronJobTestCase{
		offset: "16:10",
		c:      c,
		out:    "19:00 today - /bin/run_me_sixty_times",
		err:    nil,
	}
}

func TestParseSchedule(t *testing.T) {
	content, err := ioutil.ReadFile("../testdata/input.txt")
	require.NoError(t, err)
	testData := string(content)

	result := ParseSchedule(testData)
	assert.Len(t, result, 4)
}
