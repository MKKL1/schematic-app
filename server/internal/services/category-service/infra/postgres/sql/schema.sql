CREATE TABLE categories (
    name TEXT PRIMARY KEY,
    value_definitions JSONB NOT NULL
);

CREATE TABLE post_category_values (
    post_id BIGINT NOT NULL,
    category TEXT NOT NULL REFERENCES categories(name) ON DELETE CASCADE,
    values JSONB NOT NULL,
    PRIMARY KEY (post_id, category)
);

CREATE INDEX idx_all_values_search ON post_category_values
    USING GIN (values jsonb_path_ops);

CREATE INDEX idx_category_lookup ON post_category_values (category);