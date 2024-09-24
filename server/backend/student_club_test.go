package backend

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCampusServer_ListStudentClub(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	exampleClubs := genExampleClubData(db, t)
	server := CampusServer{db: db}
	language := pb.Language_German
	response, err := server.ListStudentClub(ctx, &pb.ListStudentClubRequest{Language: &language})
	require.NoError(t, err)
	url0 := exampleClubs[0].Image.FullExternalUrl()
	expectedResp := &pb.ListStudentClubReply{
		Collections: []*pb.StudentClubCollection{
			{
				UnstableCollectionId: uint64(exampleClubs[0].StudentClubCollection.ID),
				Title:                exampleClubs[0].StudentClubCollection.Name,
				Description:          exampleClubs[0].StudentClubCollection.Description,
				Clubs: []*pb.StudentClub{
					{
						Name:        exampleClubs[0].Name,
						Description: exampleClubs[0].Description.Ptr(),
						LinkUrl:     exampleClubs[0].LinkUrl.Ptr(),
						CoverUrl:    &url0,
					},
					{
						Name: exampleClubs[1].Name,
					},
				},
			},
			{
				UnstableCollectionId: uint64(exampleClubs[2].StudentClubCollection.ID),
				Title:                exampleClubs[2].StudentClubCollection.Name,
				Description:          exampleClubs[2].StudentClubCollection.Description,
				Clubs: []*pb.StudentClub{
					{
						Name: exampleClubs[2].Name,
					},
				},
			},
		},
	}
	require.Equal(t, expectedResp, response)
}

func genExampleClubData(db *gorm.DB, t *testing.T) []*model.StudentClub {
	col1 := model.StudentClubCollection{
		Name:        "col1",
		Language:    pb.Language_German.String(),
		Description: "Awesome collection",
	}
	err := db.Create(&col1).Error
	require.NoError(t, err)
	col2 := model.StudentClubCollection{
		Name:        "col2",
		Description: "Terrible collection",
		Language:    pb.Language_German.String(),
	}
	err = db.Create(&col2).Error
	require.NoError(t, err)
	file1 := &model.File{
		File:       1,
		Name:       fmt.Sprintf("src_%d.png", 1),
		Path:       "student_club/",
		Downloads:  1,
		URL:        null.String{},
		Downloaded: null.BoolFrom(true),
	}
	err = db.Create(file1).Error
	require.NoError(t, err)
	club1 := &model.StudentClub{
		Name:                    "Student Club 1",
		Language:                pb.Language_German.String(),
		Description:             null.StringFrom("With Description"),
		LinkUrl:                 null.StringFrom("https://example.com"),
		ImageID:                 null.IntFrom(file1.File),
		Image:                   file1,
		ImageCaption:            null.StringFrom("source: idk, something"),
		StudentClubCollectionID: col1.ID,
		StudentClubCollection:   col1,
	}
	err = db.Create(club1).Error
	require.NoError(t, err)
	club2 := &model.StudentClub{
		Name:                    "Student Club 2",
		Language:                pb.Language_German.String(),
		StudentClubCollectionID: col1.ID,
		StudentClubCollection:   col1,
	}
	err = db.Create(club2).Error
	require.NoError(t, err)
	club3 := &model.StudentClub{
		Name:                    "Student Club 3",
		Language:                pb.Language_German.String(),
		StudentClubCollectionID: col2.ID,
		StudentClubCollection:   col2,
	}
	err = db.Create(club3).Error
	require.NoError(t, err)
	return []*model.StudentClub{club1, club2, club3}
}
