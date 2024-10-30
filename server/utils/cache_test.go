package utils

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"testing"
	"time"
)

// test Cache
func TestCache(t *testing.T) {
	cache := NewCache()
	cache.Set(CacheKeyAllNewsSources, "", []model.NewsSource{}, 1*time.Hour)
	cache.Set(CacheKeyNews, "1", []model.News{}, 1*time.Hour)

	if !cache.Exists(CacheKeyAllNewsSources, "") {
		t.Errorf("CacheKeyAllNewsSources should exist")
	}
	if !cache.Exists(CacheKeyNews, "1") {
		t.Errorf("CacheKeyNews should exist")
	}

	if cache.Get(CacheKeyAllNewsSources, "").([]model.NewsSource) == nil {
		t.Errorf("CacheKeyAllNewsSources should not be nil")
	}
	if cache.Get(CacheKeyNews, "1").([]model.News) == nil {
		t.Errorf("CacheKeyNews should not be nil")
	}

	cache.delete(string(CacheKeyAllNewsSources))

	if cache.Exists(CacheKeyAllNewsSources, "") {
		t.Errorf("CacheKeyAllNewsSources should not exist")
	}
}
