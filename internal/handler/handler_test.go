package handler

import (
	"errors"
	"testing"

	"github.com/TestardR/cron_scheduler/internal/cron"
	"github.com/TestardR/cron_scheduler/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type handlerTestCase struct {
	h   Handler
	err error
}

func TestHandlerStart(t *testing.T) {
	mc := gomock.NewController(t)

	tests := map[string]handlerTestCase{
		"success":        handlerCaseStartSuccess(mc),
		"schedule-error": handlerCaseScheduleFailed(mc),
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			err := tc.h.Start()
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func handlerCaseStartSuccess(mc *gomock.Controller) handlerTestCase {
	mLogger := mock.NewMockLogger(mc)
	mScheduler := mock.NewMockScheduler(mc)

	mOffset := "10:10"

	mLogger.EXPECT().Info("scheduling cronjobs started")
	mScheduler.EXPECT().Schedule(mOffset).Return("test", nil)
	mLogger.EXPECT().Info("scheduling cronjobs finished")

	return handlerTestCase{
		h: Handler{
			log: mLogger,
			cronjobs: []cron.Scheduler{
				mScheduler,
			},
			offset: mOffset,
		},
		err: nil,
	}
}

func handlerCaseScheduleFailed(mc *gomock.Controller) handlerTestCase {
	mLogger := mock.NewMockLogger(mc)
	mScheduler := mock.NewMockScheduler(mc)

	mOffset := "10:10"
	mError := errors.New("mock")

	mLogger.EXPECT().Info("scheduling cronjobs started")
	mScheduler.EXPECT().Schedule(mOffset).Return("", mError)
	mLogger.EXPECT().Error(gomock.Any())

	return handlerTestCase{
		h: Handler{
			log: mLogger,
			cronjobs: []cron.Scheduler{
				mScheduler,
			},
			offset: mOffset,
		},
		err: errFailedCronjob,
	}
}
