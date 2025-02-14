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
WITH ins_post AS (
    INSERT INTO post (id, name, "desc", owner, author_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
),
     ins_tags AS (
         INSERT INTO post_tags (post_id, tag)
             SELECT ins_post.id, t
             FROM ins_post, unnest($6::text[]) AS t
     ),
     ins_cat AS (
         INSERT INTO post_category_metadata (post_id, category, metadata)
             SELECT ins_post.id, r."Name", r."Metadata"
             FROM ins_post,
                  jsonb_to_recordset($7::jsonb) AS r("Name" text, "Metadata" jsonb)
     )
SELECT id FROM ins_post;

-- name: GetCategory :one
SELECT * FROM categories
WHERE name = $1;