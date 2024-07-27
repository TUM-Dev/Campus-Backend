package cron

import (
	"context"
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
	db.WithContext(ctx).Find(&clubs)
	require.Equal(t, []model.StudentClub{}, clubs)
	require.NoError(t, service.studentClubCron())
	db.WithContext(ctx).Find(&clubs)
	require.True(t, len(clubs) > 0)
}
