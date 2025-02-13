-- name: GetPost :one
SELECT
    p.id,
    p.name,
    p."desc" AS description,
    p.owner,
    p.author_id,
    COALESCE(
            (
                SELECT json_agg(
                               json_build_object(
                                       'name', pcm.category,
                                       'metadata', pcm.metadata
                               )
                       )::text
                FROM post_category_metadata pcm
                WHERE pcm.post_id = p.id
            ),
            '[]'
    ) AS category_vars,
    COALESCE(
            (
                SELECT array_agg(pt.tag)::text[]
                FROM post_tags pt
                WHERE pt.post_id = p.id
            ),
            '{}'::text[]
    ) AS tags
FROM post p
WHERE p.id = $1;


-- name: CreatePost :exec
WITH new_post AS (
    INSERT INTO post (name, "desc", owner, author_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id
),
insert_tags AS (
    INSERT INTO post_tags (post_id, tag)
    SELECT np.id, t
    FROM new_post np, unnest($5::text[]) AS t
),
insert_category AS (
    INSERT INTO post_category_metadata (post_id, category, metadata)
    SELECT np.id, pair.category, pair.metadata
    FROM new_post np, unnest($6::category_metadata_pair[]) AS pair
)
SELECT * FROM new_post;