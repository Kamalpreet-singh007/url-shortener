
import(
	"database/sql"
	"context"
)

type PostgresStore struct{
	db *sql.DB
}

func (s *PostgresStore)GetByShortCode(ctx context.Context, shortcode string)(*URL, error){
	row:= s.db.QueryRowContext(ctx, "SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1",shortcode)

	var u URL
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil  // not found, not an error
	}
	if err != nil {
		return nil, err  // real problem
	}
	return &u, nil  // happy path
}

func(s * PostgresStore) CreateUrl(ctx context.Context, url string, shortcode string)(string , error){
	row := s.db.QueryRowContext(ctx,
		"INSERT INTO urls (original_url, short_code) VALUES ($1, $2) RETURNING id, original_url, short_code, created_at",
		url,shortcode
	)
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}