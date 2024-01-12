package main

import (
	"os"

	"github.com/Kamila3820/go-shop-tutorial/config"
	"github.com/Kamila3820/go-shop-tutorial/modules/servers"
	"github.com/Kamila3820/go-shop-tutorial/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
