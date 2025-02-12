-- name: GetTagsForPost :many
SELECT tag FROM post_tags
WHERE post_id = $1;

-- name: CountPostsForTag :one
SELECT count(*) FROM post_tags
WHERE tag = $1;