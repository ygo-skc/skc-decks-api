package main

import (
	"github.com/ygo-skc/skc-deck-api/api"
	"github.com/ygo-skc/skc-deck-api/db"
)

func main() {
	db.EstablishSKCDeckAPIDBConn()
	go api.RunHttpServer()
	select {}
}
