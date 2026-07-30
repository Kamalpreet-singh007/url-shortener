package store

import(
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type CacheStore struct{
	dbStore URLStore
	cache Cache	
}

func NewCacheStore(dbStore URLStore, cache Cache) *CacheStore {
	return &CacheStore{
		dbStore: dbStore,
		cache: cache,
	}
}

func (cs *CacheStore) GetByShortCode(ctx context.Context, shortcode string)(*URL, error){
	shortcodeKey := "url_shortener:shortcode:" + shortcode
	url,err := cs.cache.get(ctx, shortcodeKey)
	if err == nil && url != "" {
		SetCacheStatus(ctx, "hit")
		return &URL{ShortCode: shortcode, OriginalURL: url}, nil
	}
	if err != nil && err != redis.Nil {
		return nil, err
	}

	SetCacheStatus(ctx, "miss")
	urlObj,err := cs.dbStore.GetByShortCode(ctx, shortcode)
	if err != nil {
		return nil, err
	}

	err = cs.cache.set(ctx, shortcodeKey, urlObj.OriginalURL, 24*time.Hour)
	if err != nil {
		return nil, err
	}	
	return urlObj, nil
}

func (cs *CacheStore) CreateUrl(ctx context.Context, url string)(*URL , error){
	urlObj,err := cs.dbStore.CreateUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	shortcodeKey := "url_shortener:shortcode:" + urlObj.ShortCode
	err = cs.cache.set(ctx, shortcodeKey, urlObj.OriginalURL, 24*time.Hour)
	if err != nil {
		return nil, err
	}
	return urlObj, nil
}