package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ygo-skc/skc-deck-api/model"
	"github.com/ygo-skc/skc-deck-api/util"
)

// Handler for status/health check endpoint of the api.
// Will get status of downstream services as well to help isolate problems.
func getAPIStatusHandler(res http.ResponseWriter, req *http.Request) {
	logger, ctx := util.NewRequestSetup(context.Background(), "status")

	downstreamHealth := []model.DownstreamItem{}

	var skcDeckAPIDB string

	// get SKC Deck API status by checking the version number. If this operation fails, its save to assume the DB is down.
	if dbVersion, err := skcDeckAPIDBInterface.GetSKCDeckAPIDBVersion(ctx); err != nil {
		downstreamHealth = append(downstreamHealth, model.DownstreamItem{ServiceName: "SKC Deck API DB", Status: model.Down})
	} else {
		downstreamHealth = append(downstreamHealth, model.DownstreamItem{ServiceName: "SKC Deck API DB", Status: model.Up})
		skcDeckAPIDB = dbVersion
	}

	status := model.APIHealth{Version: "1.0.0", Downstream: downstreamHealth}

	logger.Info(fmt.Sprintf("API Status Info! SKC Deck API DB version: %s", skcDeckAPIDB))
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(status)
}
