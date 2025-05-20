-- Таблица стран
CREATE TABLE IF NOT EXISTS countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capital VARCHAR(100) NOT NULL,
    language VARCHAR(50),
    currency VARCHAR(50),
    description TEXT,
    photo_url VARCHAR(255)
);

-- Таблица мест
CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    longitude DECIMAL(10, 6),
    latitude DECIMAL(10, 6)
);

-- Таблица фотографий мест (для хранения массива URL)
CREATE TABLE IF NOT EXISTS place_photos (
    id SERIAL PRIMARY KEY,
    place_id INTEGER NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_countries_name ON countries(name);
CREATE INDEX IF NOT EXISTS idx_places_name ON places(name);
CREATE INDEX IF NOT EXISTS idx_place_photos_place_id ON place_photos(place_id);