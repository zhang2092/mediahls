package db

import "context"

type CreateVideoTxParam struct {
	CreateVideoParams
	AfterCreate func(video Video) error
}

type CreateVideoTxResult struct {
	Video Video
}

func (store *SQLStore) CreateVideoTx(ctx context.Context, arg CreateVideoTxParam) (CreateVideoTxResult, error) {
	var result CreateVideoTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		result.Video, err = q.CreateVideo(ctx, arg.CreateVideoParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Video)
	})

	return result, err
}
