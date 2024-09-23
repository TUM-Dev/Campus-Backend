package cron

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"testing"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/stretchr/testify/require"
)

func TestCronService_studentClubCron(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	service := New(db)
	var clubs []model.StudentClub
	require.NoError(t, db.WithContext(ctx).Find(&clubs).Error)
	require.Equal(t, []model.StudentClub{}, clubs)
	require.NoError(t, service.studentClubCron(pb.Language_German))
	require.NoError(t, db.WithContext(ctx).Find(&clubs).Error)
	require.True(t, len(clubs) > 0)
	require.NoError(t, db.WithContext(ctx).Where("1=1").Delete(&model.StudentClub{}).Error)
	require.NoError(t, service.studentClubCron(pb.Language_English))
	require.NoError(t, db.WithContext(ctx).Find(&clubs).Error)
	require.True(t, len(clubs) > 0)
}
