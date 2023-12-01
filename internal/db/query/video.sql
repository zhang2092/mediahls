-- name: CreateVideo :one
INSERT INTO videos (
  id, title, description, images, origin_link, play_link, user_id, create_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: DeleteVideo :exec
UPDATE videos
SET is_deleted = TRUE
WHERE id = $1;

-- name: UpdateVideoStatus :one
UPDATE videos
SET status = $2,
    update_at = $3,
    update_by = $4
WHERE id = $1
RETURNING *;

-- name: SetVideoPlay :one
UPDATE videos
SET status = $2,
    play_link = $3,
    update_at = $4,
    update_by = $5
WHERE id = $1
RETURNING *;

-- name: UpdateVideo :one
UPDATE videos
SET title = $2,
    description = $3,
    images = $4,
    status = $5,
    update_at = $6,
    update_by = $7
WHERE id = $1
RETURNING *;

-- name: GetVideo :one
SELECT * FROM videos
WHERE id = $1 LIMIT 1;

-- name: ListVideos :many
SELECT * FROM videos
WHERE is_deleted = FALSE AND status=200
ORDER BY id DESC
LIMIT $1
OFFSET $2;

-- name: ListVideosByUser :many
SELECT * FROM videos
WHERE is_deleted = FALSE AND user_id = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;