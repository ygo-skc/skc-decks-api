package downstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cModel "github.com/ygo-skc/skc-go/common/model"
	cUtil "github.com/ygo-skc/skc-go/common/util"
)

const (
	BATCH_CARD_INFO_ENDPOINT  = "/api/v1/suggestions/card-details"
	BATCH_CARD_INFO_OPERATION = "Batch Card Info"
	BATCH_CARD_INFO_ERROR     = "There was an error fetching card info"
)

func FetchBatchCardData(ctx context.Context, cardIDs []string) (*cModel.BatchCardData[cModel.CardIDs], *cModel.APIError) {
	logger := cUtil.LoggerFromContext(ctx)
	logger.Info(fmt.Sprintf("Fetching card info for the following IDs: %v", cardIDs))

	var resp *http.Response
	var err error

	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(cModel.BatchCardIDs{CardIDs: cardIDs})

	if resp, err = suggestionEngineClient.Post(
		fmt.Sprintf("https://skc-suggestion-engine:9000%s", BATCH_CARD_INFO_ENDPOINT), "application/json", reqBody); err != nil {
		logger.Error(
			fmt.Sprintf("There was an issue calling Suggestion Engine. Operation: %s. Error: %s",
				BATCH_CARD_INFO_OPERATION, err))
		return nil, &cModel.APIError{Message: BATCH_CARD_INFO_ERROR, StatusCode: http.StatusInternalServerError}
	} else {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		s, _ := io.ReadAll(resp.Body)
		logger.Error(
			fmt.Sprintf("Suggestion Engine returned with non 200 status. Operation: %s. Body: %s. Code: %d",
				BATCH_CARD_INFO_OPERATION, string(s), resp.StatusCode))
		return nil, &cModel.APIError{Message: BATCH_CARD_INFO_ERROR, StatusCode: http.StatusInternalServerError}
	}

	var cardData cModel.BatchCardData[cModel.CardIDs]
	if err = json.NewDecoder(resp.Body).Decode(&cardData); err != nil && err != io.EOF {
		logger.Error(
			fmt.Sprintf("Error occurred while deserializing output from Suggestion Engine. Operation: %s. Error %v",
				BATCH_CARD_INFO_OPERATION, err))
		return nil, &cModel.APIError{Message: BATCH_CARD_INFO_ERROR, StatusCode: http.StatusInternalServerError}
	}

	return &cardData, nil
}
