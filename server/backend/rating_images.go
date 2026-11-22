package backend

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const ratingImageBaseURL = "https://api.tum.app/rating-images"

func dishRatingImageURL(ratingID int64) string {
	if ratingID <= 0 {
		return ""
	}
	return fmt.Sprintf("%s/dish/%d", ratingImageBaseURL, ratingID)
}

func canteenRatingImageURL(ratingID int64) string {
	if ratingID <= 0 {
		return ""
	}
	return fmt.Sprintf("%s/canteen/%d", ratingImageBaseURL, ratingID)
}

// RatingImageHandler serves rating photos at /rating-images/{dish|canteen}/{rating_id}.
// It only serves images belonging to visible ratings, keeping moderated/deleted photos private
// while /files/rating/ remains blocked by PublicFileServer.
func (s *CampusServer) RatingImageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if s.db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		trimmed := strings.TrimPrefix(r.URL.Path, "/")
		kind, idStr, ok := strings.Cut(trimmed, "/")
		if !ok || strings.Contains(idStr, "/") {
			http.NotFound(w, r)
			return
		}
		ratingID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || ratingID <= 0 {
			http.NotFound(w, r)
			return
		}

		var imagePath string
		db := s.db.WithContext(r.Context())
		switch kind {
		case "dish":
			var rating model.DishRating
			err = db.Select("image").First(&rating, "dishRating = ? AND status = ?", ratingID, model.RatingStatusVisible).Error
			imagePath = rating.Image
		case "canteen":
			var rating model.CanteenRating
			err = db.Select("image").First(&rating, "cafeteriaRating = ? AND status = ?", ratingID, model.RatingStatusVisible).Error
			imagePath = rating.Image
		default:
			http.NotFound(w, r)
			return
		}
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
				return
			}
			log.WithError(err).WithFields(log.Fields{"kind": kind, "ratingID": ratingID}).Error("failed to look up rating image")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if imagePath == "" {
			http.NotFound(w, r)
			return
		}

		data, err := os.ReadFile(imagePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				http.NotFound(w, r)
				return
			}
			log.WithError(err).WithField("path", imagePath).Error("failed to read rating image")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		sum := sha256.Sum256(data)
		etag := `"` + hex.EncodeToString(sum[:16]) + `"`
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodGet {
			_, _ = w.Write(data)
		}
	})
}
