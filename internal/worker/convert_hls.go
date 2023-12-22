package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/pkg/convert"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
)

const TaskConvertHLS = "task:convert_hls"

type PayloadConvertHLS struct {
	Id string `json:"id"`
}

func (distributor *RedisTaskDistributor) DistributeConvertHLS(
	ctx context.Context,
	payload *PayloadConvertHLS,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskConvertHLS, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("type: %s\n", task.Type())
	log.Printf("payload: %s\n", task.Payload())
	log.Printf("queue: %s\n", info.Queue)
	log.Printf("max_retry: %d\n", info.MaxRetry)
	log.Printf("enqueued task\n")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskConvertHLS(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadConvertHLS
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	video, err := processor.store.GetVideo(ctx, payload.Id)
	if err != nil {
		return fmt.Errorf("failed to get video by id [%s] in db: %w", payload.Id, err)
	}

	arg := db.UpdateVideoStatusParams{
		ID:       video.ID,
		Status:   1,
		UpdateAt: time.Now(),
		UpdateBy: "任务队列",
	}
	video, err = processor.store.UpdateVideoStatus(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to video by id [%s]: in db: %w", payload.Id, err)
	}

	err = convert.ConvertHLS("media/"+video.ID+"/", strings.TrimPrefix(video.OriginLink, "/"))
	if err != nil {
		logger.Logger.Errorf("Convert HLS [%s]-[%s]: %v", video.ID, video.OriginLink, err)
		arg = db.UpdateVideoStatusParams{
			ID:       video.ID,
			Status:   2,
			UpdateAt: time.Now(),
			UpdateBy: "任务队列",
		}
		_, _ = processor.store.UpdateVideoStatus(ctx, arg)
		return fmt.Errorf("failed to convert hls by [%s]: %w", payload.Id, err)
	}

	// 转码成功
	if _, err = processor.store.SetVideoPlay(ctx, db.SetVideoPlayParams{
		ID:       video.ID,
		Status:   200,
		PlayLink: "/media/" + video.ID + "/stream/",
		UpdateAt: time.Now(),
		UpdateBy: "任务队列",
	}); err != nil {
		logger.Logger.Errorf("Set Video Play [%s]-[%s]: %v", video.ID, video.OriginLink, err)
		return fmt.Errorf("failed to set video [%s] play: %w", video.ID, err)
	}

	logger.Logger.Infof("[%s]-[%s] 转码完成", video.ID, video.OriginLink)
	return nil
}
