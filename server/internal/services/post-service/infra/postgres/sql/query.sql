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
    ) AS tags,
    COALESCE(
            (
                SELECT json_agg(
                               json_build_object(
                                       'hash', af.hash,
                                       'temp_id', af.temp_id,
                                       'name', af.name,
                                       'file_size', af.file_size,
                                       'downloads', af.downloads,
                                       'created_at', af.created_at,
                                       'updated_at', af.updated_at
                               )
                       )::text
                FROM attached_files af
                WHERE af.post_id = p.id
            ),
            '[]'
    ) AS files
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
     ),
     ins_file AS (
         INSERT INTO attached_files (temp_id, post_id, name)
             SELECT f.temp_id, ins_post.id, f.name
             FROM ins_post,
                  jsonb_to_recordset($8::jsonb) AS f(temp_id uuid, name text)
     )
SELECT id FROM ins_post;

-- name: GetCategory :one
SELECT * FROM categories
WHERE name = $1;