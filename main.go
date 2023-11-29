package main

import (
	"database/sql"
	"log"

	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/handlers"
	"github.com/zhang2092/mediahls/internal/pkg/config"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
)

func main() {
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
	server, err := handlers.NewServer(config, store)
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
