CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    lecturer VARCHAR(100) NOT NULL,
    capacity INT NOT NULL CHECK (capacity >= 0),
    taken INT NOT NULL DEFAULT 0 CHECK (taken >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT taken_tidak_melebihi_kapasitas CHECK (taken <= capacity)
);