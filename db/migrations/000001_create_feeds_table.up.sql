CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    item_selector TEXT,
    title_selector TEXT,
    description_selector TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_up TIMESTAMP DEFAULT NULL
);
