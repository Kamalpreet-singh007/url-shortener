package store

import(
	"context"
	"time"
)

type URL struct {
	ID			 int64
	OriginalURL string
	ShortCode 	 string
	CreatedAt 	 time.Time
	
}

type URLStore interface {
	CreateUrl(ctx context.Context, originalurl string)(*URL , error)
	GetByShortCode(ctx context.Context, shortcode string)(*URL, error)

}
