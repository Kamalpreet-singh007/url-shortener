package store

import(
	"database/sql"
	"context"
	"github.com/Kamalpreet-singh007/url-shortener/pkg/base62"
)

var _ URLStore = (*PostgresStore)(nil)

type PostgresStore struct{
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
    return &PostgresStore{db: db}
}


func (s *PostgresStore)GetByShortCode(ctx context.Context, shortcode string)(*URL, error){
	row:= s.db.QueryRowContext(ctx, "SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1",shortcode)

	var u URL
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt)
	if err != nil {
		return nil, err  // real problem
	}
	return &u, nil  // happy path
}

func(s * PostgresStore) CreateUrl(ctx context.Context, url string)(*URL , error){
	var id int64
	err := s.db.QueryRowContext(ctx, "SELECT nextval('urls_id_seq')").Scan(&id)
	if err != nil {
		return nil, err
	}
	shortcode := base62.Encode(id)

	row := s.db.QueryRowContext(ctx,`INSERT INTO urls (id, original_url, short_code) 
									 VALUES ($1, $2, $3)
									 on conflict(original_url) do update
									 set original_url = excluded.original_url
									 RETURNING id, original_url, short_code, created_at`, id, url, shortcode)

	var u URL
	err = row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}