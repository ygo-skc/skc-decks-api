package model

import (
	"strings"
)

type BatchCardIDs struct {
	CardIDs []string `json:"cardIDs" validate:"required"`
}

type BatchCardInfo struct {
	CardInfo       CardDataMap `json:"cardInfo"`
	InvalidCardIDs []string    `json:"invalidCardIDs"`
}

type Card struct {
	CardID         string  `db:"card_number" json:"cardID"`
	CardColor      string  `db:"card_color" json:"cardColor"`
	CardName       string  `db:"card_name" json:"cardName"`
	CardAttribute  string  `db:"card_attribute" json:"cardAttribute"`
	CardEffect     string  `db:"card_effect" json:"cardEffect"`
	MonsterType    *string `db:"monster_type" json:"monsterType,omitempty"`
	MonsterAttack  *uint16 `db:"monster_attack" json:"monsterAttack,omitempty"`
	MonsterDefense *uint16 `db:"monster_defense" json:"monsterDefense,omitempty"`
}

type CardDataMap map[string]Card

func (c Card) IsExtraDeckMonster() bool {
	color := strings.ToUpper(c.CardColor)
	return strings.Contains(color, "FUSION") || strings.Contains(color, "SYNCHRO") || strings.Contains(color, "XYZ") || strings.Contains(color, "PENDULUM") || strings.Contains(color, "LINK")
}
