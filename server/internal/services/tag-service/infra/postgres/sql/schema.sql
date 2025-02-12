CREATE TABLE post_tags (
    post_id BIGINT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (post_id, tag)
);

CREATE INDEX idx_post_tags_tag ON post_tags(tag);