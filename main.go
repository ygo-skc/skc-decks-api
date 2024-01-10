package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ygo-skc/skc-deck-api/api"
	"github.com/ygo-skc/skc-deck-api/db"
	"github.com/ygo-skc/skc-deck-api/util"
)

func main() {
	util.SetupEnv()
	db.EstablishSKCDeckAPIDBConn()

	api.ConfigureServer()
	go api.ServeTLS()
	api.ServeUnsecured()
}
