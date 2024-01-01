package main

import (
	"github.com/ygo-skc/skc-deck-api/api"
)

func main() {
	api.ConfigureServer()
	go api.ServeTLS()
	api.ServeUnsecured()
}
