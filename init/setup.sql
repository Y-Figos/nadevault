CREATE TABLE IF NOT EXISTS maps (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL
);

INSERT INTO maps (slug, name) VALUES ('mirage', 'Mirage');