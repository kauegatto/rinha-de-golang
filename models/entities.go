package models

import (
	"errors"
	"fmt"
	"time"
)

type Payment struct {
	Valor     int64  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
	CriadoEm  time.Time
}

func (p Payment) Validar(saldo int64) (valor int64, error error) {
	if p.Tipo != "c" && p.Tipo != "d" {
		return 0, fmt.Errorf("invalid type %s should be either p or s", p.Tipo)
	}

	if len(p.Descricao) > 10 {
		return 0, errors.New("description should not be longer than 10 characters")
	}

	newBalance := p.Valor - saldo
	if newBalance < 0 {
		return 0, fmt.Errorf("would let inconsistant amount (negative): %d", p.Valor)
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
