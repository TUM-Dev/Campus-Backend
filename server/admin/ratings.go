package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
)

type ratingTag struct {
	RatingID int64
	Label    string
	Points   int32
}

type ratingEntry struct {
	ID        int64
	Canteen   string
	Points    int32
	Comment   string
	Timestamp time.Time
	Tags      []ratingTag `gorm:"-"`
}

type ratingsPage struct {
	page
	pagination
	Canteens   []model.Canteen
	SelectedID int64
	Count      int64
	Average    null.Float
	Entries    []ratingEntry
}

func (s *server) ratings(w http.ResponseWriter, r *http.Request, session model.AdminSession) {
	paging, err := requestedPage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := s.db.WithContext(r.Context())
	data := ratingsPage{page: basePage("Cafeteria ratings", "ratings", session), pagination: paging}
	if err := db.Order("name").Find(&data.Canteens).Error; err != nil {
		s.internalError(w, err)
		return
	}
	if value := r.URL.Query().Get("canteen"); value != "" {
		id, err := strconv.ParseInt(value, 10, 32)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid cafeteria ID.", http.StatusBadRequest)
			return
		}
		found := false
		for _, canteen := range data.Canteens {
			if canteen.CanteenID == id {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "Cafeteria not found.", http.StatusNotFound)
			return
		}
		data.SelectedID = id
	}
	query := db.Model(&model.CanteenRating{})
	if data.SelectedID != 0 {
		query = query.Where("canteen_ratings.cafeteriaID = ?", data.SelectedID)
	}
	var stats struct {
		Count   int64
		Average null.Float
	}
	if err := query.Select("COUNT(*) AS count, AVG(points) AS average").Scan(&stats).Error; err != nil {
		s.internalError(w, err)
		return
	}
	data.Count, data.Average = stats.Count, stats.Average
	if err := query.Select("canteen_ratings.cafeteriaRating AS id, canteen_ratings.points, canteen_ratings.comment, canteen_ratings.timestamp, canteens.name AS canteen").
		Joins("JOIN canteens ON canteens.cafeteria = canteen_ratings.cafeteriaID").
		Order("canteen_ratings.timestamp DESC").Order("canteen_ratings.cafeteriaRating DESC").
		Offset((paging.Number - 1) * pageSize).Limit(pageSize + 1).Scan(&data.Entries).Error; err != nil {
		s.internalError(w, err)
		return
	}
	data.HasNext = len(data.Entries) > pageSize
	if data.HasNext {
		data.Entries = data.Entries[:pageSize]
	}
	ids := make([]int64, 0, len(data.Entries))
	for _, entry := range data.Entries {
		ids = append(ids, entry.ID)
	}
	if len(ids) > 0 {
		var tags []ratingTag
		if err := db.Table("canteen_rating_tags AS tags").
			Joins("JOIN canteen_rating_tag_options AS options ON options.cafeteriaRatingTagOption = tags.tagID").
			Select("tags.correspondingRating AS rating_id, options.EN AS label, tags.points").
			Where("tags.correspondingRating IN ?", ids).Order("options.EN").Scan(&tags).Error; err != nil {
			s.internalError(w, err)
			return
		}
		for i := range data.Entries {
			for _, tag := range tags {
				if tag.RatingID == data.Entries[i].ID {
					data.Entries[i].Tags = append(data.Entries[i].Tags, tag)
				}
			}
		}
	}
	s.render(w, "ratings", data)
}
