package models

import (
	"errors"
	"fmt"
	"time"
)

type Payment struct {
	Value       int64  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
	CreatedAt   time.Time
}

func (p Payment) ValidateAndReturnNewBalance(balance int64) (valor int64, error error) {
	if p.Type != "c" && p.Type != "d" {
		return 0, fmt.Errorf("invalid type %s should be either p or s", p.Type)
	}

	if len(p.Description) > 10 {
		return 0, errors.New("description should not be longer than 10 characters")
	}

	if p.Type == "d" {
		p.Value = p.Value * -1
	}

	newBalance := balance + p.Value

	if newBalance < 0 {
		return 0, fmt.Errorf("would let inconsistant amount (negative): %d", p.Value)
	}

	return newBalance, nil
}

type Client struct {
	Date         string `json:"data_extrato"`
	AccountLimit int64  `json:"limite"`
	Balance      int64  `json:"total"`
}

type Transaction struct {
	Amount      int64     `json:"valor"`
	Operation   string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}
