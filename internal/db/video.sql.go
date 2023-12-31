// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: video.sql

package db

import (
	"context"
	"time"
)

const createVideo = `-- name: CreateVideo :one
INSERT INTO videos (
  id, title, description, images, origin_link, play_link, user_id, create_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by
`

type CreateVideoParams struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Images      string `json:"images"`
	OriginLink  string `json:"origin_link"`
	PlayLink    string `json:"play_link"`
	UserID      string `json:"user_id"`
	CreateBy    string `json:"create_by"`
}

func (q *Queries) CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error) {
	row := q.db.QueryRowContext(ctx, createVideo,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Images,
		arg.OriginLink,
		arg.PlayLink,
		arg.UserID,
		arg.CreateBy,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Images,
		&i.OriginLink,
		&i.PlayLink,
		&i.Status,
		&i.IsDeleted,
		&i.UserID,
		&i.CreateAt,
		&i.CreateBy,
		&i.UpdateAt,
		&i.UpdateBy,
	)
	return i, err
}

const deleteVideo = `-- name: DeleteVideo :exec
UPDATE videos
SET is_deleted = TRUE
WHERE id = $1
`

func (q *Queries) DeleteVideo(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteVideo, id)
	return err
}

const getVideo = `-- name: GetVideo :one
SELECT id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by FROM videos
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetVideo(ctx context.Context, id string) (Video, error) {
	row := q.db.QueryRowContext(ctx, getVideo, id)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Images,
		&i.OriginLink,
		&i.PlayLink,
		&i.Status,
		&i.IsDeleted,
		&i.UserID,
		&i.CreateAt,
		&i.CreateBy,
		&i.UpdateAt,
		&i.UpdateBy,
	)
	return i, err
}

const listVideos = `-- name: ListVideos :many
SELECT id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by FROM videos
WHERE is_deleted = FALSE AND status=200
ORDER BY id DESC
LIMIT $1
OFFSET $2
`

type ListVideosParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListVideos(ctx context.Context, arg ListVideosParams) ([]Video, error) {
	rows, err := q.db.QueryContext(ctx, listVideos, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Video{}
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Images,
			&i.OriginLink,
			&i.PlayLink,
			&i.Status,
			&i.IsDeleted,
			&i.UserID,
			&i.CreateAt,
			&i.CreateBy,
			&i.UpdateAt,
			&i.UpdateBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVideosByUser = `-- name: ListVideosByUser :many
SELECT id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by FROM videos
WHERE is_deleted = FALSE AND user_id = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3
`

type ListVideosByUserParams struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListVideosByUser(ctx context.Context, arg ListVideosByUserParams) ([]Video, error) {
	rows, err := q.db.QueryContext(ctx, listVideosByUser, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Video{}
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Images,
			&i.OriginLink,
			&i.PlayLink,
			&i.Status,
			&i.IsDeleted,
			&i.UserID,
			&i.CreateAt,
			&i.CreateBy,
			&i.UpdateAt,
			&i.UpdateBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setVideoPlay = `-- name: SetVideoPlay :one
UPDATE videos
SET status = $2,
    play_link = $3,
    update_at = $4,
    update_by = $5
WHERE id = $1
RETURNING id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by
`

type SetVideoPlayParams struct {
	ID       string    `json:"id"`
	Status   int32     `json:"status"`
	PlayLink string    `json:"play_link"`
	UpdateAt time.Time `json:"update_at"`
	UpdateBy string    `json:"update_by"`
}

func (q *Queries) SetVideoPlay(ctx context.Context, arg SetVideoPlayParams) (Video, error) {
	row := q.db.QueryRowContext(ctx, setVideoPlay,
		arg.ID,
		arg.Status,
		arg.PlayLink,
		arg.UpdateAt,
		arg.UpdateBy,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Images,
		&i.OriginLink,
		&i.PlayLink,
		&i.Status,
		&i.IsDeleted,
		&i.UserID,
		&i.CreateAt,
		&i.CreateBy,
		&i.UpdateAt,
		&i.UpdateBy,
	)
	return i, err
}

const updateVideo = `-- name: UpdateVideo :one
UPDATE videos
SET title = $2,
    description = $3,
    images = $4,
    status = $5,
    update_at = $6,
    update_by = $7
WHERE id = $1
RETURNING id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by
`

type UpdateVideoParams struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	Status      int32     `json:"status"`
	UpdateAt    time.Time `json:"update_at"`
	UpdateBy    string    `json:"update_by"`
}

func (q *Queries) UpdateVideo(ctx context.Context, arg UpdateVideoParams) (Video, error) {
	row := q.db.QueryRowContext(ctx, updateVideo,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Images,
		arg.Status,
		arg.UpdateAt,
		arg.UpdateBy,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Images,
		&i.OriginLink,
		&i.PlayLink,
		&i.Status,
		&i.IsDeleted,
		&i.UserID,
		&i.CreateAt,
		&i.CreateBy,
		&i.UpdateAt,
		&i.UpdateBy,
	)
	return i, err
}

const updateVideoStatus = `-- name: UpdateVideoStatus :one
UPDATE videos
SET status = $2,
    update_at = $3,
    update_by = $4
WHERE id = $1
RETURNING id, title, description, images, origin_link, play_link, status, is_deleted, user_id, create_at, create_by, update_at, update_by
`

type UpdateVideoStatusParams struct {
	ID       string    `json:"id"`
	Status   int32     `json:"status"`
	UpdateAt time.Time `json:"update_at"`
	UpdateBy string    `json:"update_by"`
}

func (q *Queries) UpdateVideoStatus(ctx context.Context, arg UpdateVideoStatusParams) (Video, error) {
	row := q.db.QueryRowContext(ctx, updateVideoStatus,
		arg.ID,
		arg.Status,
		arg.UpdateAt,
		arg.UpdateBy,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Images,
		&i.OriginLink,
		&i.PlayLink,
		&i.Status,
		&i.IsDeleted,
		&i.UserID,
		&i.CreateAt,
		&i.CreateBy,
		&i.UpdateAt,
		&i.UpdateBy,
	)
	return i, err
}
