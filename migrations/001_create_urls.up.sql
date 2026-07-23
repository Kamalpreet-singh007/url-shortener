CREATE TABLE urls (
    ID BIGSERIAL  PRIMARY KEY,
    original_url text NOT NULL unique,
    short_code varchar(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()   
)