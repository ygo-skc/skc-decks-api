package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ygo-skc/skc-deck-api/model"
)

// Handler for status/health check endpoint of the api.
// Will get status of downstream services as well to help isolate problems.
func getAPIStatusHandler(res http.ResponseWriter, req *http.Request) {
	downstreamHealth := []model.DownstreamItem{}

	var skcDeckAPIDB string

	// get SKC Deck API status by checking the version number. If this operation fails, its save to assume the DB is down.
	if dbVersion, err := skcDeckAPIDBInterface.GetSKCDeckAPIDBVersion(); err != nil {
		downstreamHealth = append(downstreamHealth, model.DownstreamItem{ServiceName: "SKC Deck API DB", Status: model.Down})
	} else {
		downstreamHealth = append(downstreamHealth, model.DownstreamItem{ServiceName: "SKC Deck API DB", Status: model.Up})
		skcDeckAPIDB = dbVersion
	}

	status := model.APIHealth{Version: "1.0.0", Downstream: downstreamHealth}

	log.Printf("API Status Info! SKC Deck API DB version: %s", skcDeckAPIDB)
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(status)
}
