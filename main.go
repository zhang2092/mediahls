package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"

	"github.com/hibiken/asynq"
	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/handlers"
	"github.com/zhang2092/mediahls/internal/pkg/config"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
	"github.com/zhang2092/mediahls/internal/worker"
)

//go:embed web/templates
var templateFS embed.FS

//go:embed web/statics
var staticFS embed.FS

func main() {
	// Set up templates
	templates, err := fs.Sub(templateFS, "web/templates")
	if err != nil {
		log.Fatal(err)
	}

	// Set up statics
	statics, err := fs.Sub(staticFS, "web/statics")
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	logger.NewLogger()

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	redisOpt := asynq.RedisClientOpt{
		Addr:     config.RDSource,
		Password: config.RDPassowrd,
		DB:       config.RDIndex,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	server, err := handlers.NewServer(templates, statics, config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	server.Start(conn)

	// s := xid.New().String()
	// log.Println(s)
	// err := convert.ConvertHLS("media/"+s+"/", "upload/20231129/o6e6qKaMdk0VC1Ys2SHnr.mp4")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("ok")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Printf("task processor start\n")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("failed to start task processor: %w", err)
	}
}
