CREATE TABLE feeds (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    item_selector TEXT,
    title_selector TEXT,
    link_selector TEXT,
    description_selector TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);
