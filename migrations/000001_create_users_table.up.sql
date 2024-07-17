-- Таблиця cats
CREATE TABLE IF NOT EXISTS cats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    years_of_experience INT NOT NULL,
    breed TEXT NOT NULL,
    salary NUMERIC NOT NULL
);

-- Таблиця missions
CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    cat_id INT REFERENCES cats(id),
    complete BOOLEAN DEFAULT FALSE
    );

-- Таблиця targets
CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    mission_id INT REFERENCES missions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    notes TEXT,
    complete BOOLEAN DEFAULT FALSE,
    UNIQUE (mission_id, name)
    );
