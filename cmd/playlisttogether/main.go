package main

import (
	"flag"
	"fmt"
	"os"
	"playlisttogether/backend/config"
	"playlisttogether/backend/web"

	"playlisttogether/backend/database"
)

func main() {
	flag.Parse()
	conf, err := config.New()

	if err != nil {
		fmt.Println("FATAL failed loading config")
		fmt.Println(err)
		os.Exit(0)
	}

	db, err := database.NewDB(conf)
	if err != nil {
		fmt.Println("Fatal: no Connection to Database")
	}

	web.Serve(conf, db)
}
