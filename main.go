package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ygo-skc/skc-deck-api/api"
	"github.com/ygo-skc/skc-deck-api/db"
)

func main() {
	db.EstablishSKCDeckAPIDBConn()
	api.RunHttpServer()
}
