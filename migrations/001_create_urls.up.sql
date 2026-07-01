CREATE TABLE urls (
    ID BIGSERIAL  PRIMARY KEY,
    orignal_url text NOT NULL ,
    short_code varchar(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()   
)