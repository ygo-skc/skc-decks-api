package model

import (
	"sort"
	"strings"
)

type CardIDs []string

type BatchCardIDs struct {
	CardIDs CardIDs `json:"cardIDs" validate:"required"`
}

type BatchCardInfo struct {
	CardInfo       CardDataMap `json:"cardInfo"`
	InvalidCardIDs CardIDs     `json:"invalidCardIDs"`
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

func (c Card) IsExtraDeckMonster() bool {
	color := strings.ToUpper(c.CardColor)
	return strings.Contains(color, "FUSION") || strings.Contains(color, "SYNCHRO") || strings.Contains(color, "XYZ") || strings.Contains(color, "PENDULUM") || strings.Contains(color, "LINK")
}

type Cards []Card

func (cards Cards) SortCardsByName() {
	sort.SliceStable(cards, func(i, j int) bool {
		return (cards)[i].CardName < (cards)[j].CardName
	})
}

type CardDataMap map[string]Card
