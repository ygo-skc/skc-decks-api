package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ygo-skc/skc-deck-api/model"
)

var (
	skcDBConn *sql.DB
)

const (
	// errors
	genericError string = "Error occurred while querying DB."

	// queries
	queryDBVersion string = "SELECT VERSION()"
)

// interface
type SKCDatabaseAccessObject interface {
	GetSKCDBVersion() (string, error)
	FindDesiredCardInDBUsingMultipleCardIDs(cards []string) (model.CardDataMap, *model.APIError)
}

// impl
type SKCDAOImplementation struct{}

// Get version of MYSQL being used by SKC DB.
func (imp SKCDAOImplementation) GetSKCDBVersion() (string, error) {
	var version string
	if err := skcDBConn.QueryRow(queryDBVersion).Scan(&version); err != nil {
		log.Println("Error getting SKC DB version", err)
		return version, err
	}

	return version, nil
}

func (imp SKCDAOImplementation) FindDesiredCardInDBUsingMultipleCardIDs(cards []string) (model.CardDataMap, *model.APIError) {
	numCards := len(cards)
	args := make([]interface{}, numCards)
	cardData := make(map[string]model.Card, numCards)

	for index, cardId := range cards {
		args[index] = cardId
	}

	query := fmt.Sprintf("SELECT card_number, card_color, card_name, card_attribute, card_effect, monster_type, monster_attack, monster_defense FROM card_info WHERE card_number IN (%s)", variablePlaceholders(numCards))

	if rows, err := skcDBConn.Query(query, args...); err != nil {
		log.Println("Error occurred while querying SKC DB for card info using 1 or more CardIDs", err)
		return nil, &model.APIError{Message: genericError}
	} else {
		for rows.Next() {
			var card model.Card
			if err := rows.Scan(&card.CardID, &card.CardColor, &card.CardName, &card.CardAttribute, &card.CardEffect, &card.MonsterType, &card.MonsterAttack, &card.MonsterDefense); err != nil {
				log.Println("Error transforming row to Card object from SKC DB while using 1 or more CardIDs", err)
				return nil, &model.APIError{Message: "Error parsing data from DB."}
			}

			cardData[card.CardID] = card
		}
	}

	return cardData, nil
}

func variablePlaceholders(totalFields int) string {
	if totalFields == 0 {
		return ""
	} else if totalFields == 1 {
		return "?"
	} else {
		return fmt.Sprintf("?%s", strings.Repeat(", ?", totalFields-1))
	}
}
