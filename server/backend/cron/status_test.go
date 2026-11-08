package cron

import (
	"context"
	"errors"
	"testing"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
)

func TestStatusDistinguishesRunningFailuresAndSuccesses(t *testing.T) {
	service := &CronService{}
	require.True(t, service.Status().LastPoll.IsZero())
	service.recordPoll(false)
	require.False(t, service.Status().PollOK)
	service.recordJob(1, false, nil)
	require.Equal(t, "Running", service.Status().Jobs[1].State)
	service.recordJob(1, true, errors.New("upstream failed"))
	failed := service.Status().Jobs[1]
	require.Equal(t, "Failed", failed.State)
	require.True(t, failed.LastSuccess.IsZero())
	require.False(t, failed.FinishedAt.IsZero())
	service.recordJob(1, false, nil)
	service.recordJob(1, true, nil)
	success := service.Status().Jobs[1]
	require.Equal(t, "Succeeded", success.State)
	require.False(t, success.LastSuccess.IsZero())
	service.recordJob(1, false, nil)
	require.Equal(t, success.LastSuccess, service.Status().Jobs[1].LastSuccess)
	copy := service.Status()
	delete(copy.Jobs, 1)
	require.Contains(t, service.Status().Jobs, int64(1))
}

func TestSchedulerDoesNotOverwriteAdminChangesOrStartTwice(t *testing.T) {
	db := utils.SetupTestContainer(context.Background(), t)
	job := model.Crontab{Type: null.StringFrom(NewsType), Interval: 7200}
	require.NoError(t, db.Create(&job).Error)
	service := New(db)
	stale := job
	require.NoError(t, db.Model(&job).UpdateColumn("interval", 0).Error)
	claimed, err := service.claimJob(&stale)
	require.NoError(t, err)
	require.False(t, claimed)
	require.NoError(t, db.First(&job, job.Cron).Error)
	require.Zero(t, job.Interval)
	require.Zero(t, job.LastRun)
	require.NoError(t, db.Model(&job).UpdateColumn("interval", 3600).Error)
	require.NoError(t, db.First(&job, job.Cron).Error)
	stale = job
	claimed, err = service.claimJob(&job)
	require.NoError(t, err)
	require.True(t, claimed)
	claimed, err = service.claimJob(&stale)
	require.NoError(t, err)
	require.False(t, claimed)
}
