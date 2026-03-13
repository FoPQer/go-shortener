CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX idx_urls_short_url ON urls(short_url);
CREATE UNIQUE INDEX idx_urls_original_url ON urls(original_url);