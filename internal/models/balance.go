package models

import (
	"log"

	"github.com/mwolff44/pyfb-router/internal/db"
)

// BalanceStatus model
type BalanceStatus struct {
	CustomerEnabled bool
	CustomerBalance float64
	CreditLimit     float64
}

// BalanceByCustomerID gets balance information from DB
func (b BalanceStatus) BalanceByCustomerID(id string) (*BalanceStatus, error) {
	bs := &BalanceStatus{}
	pool := db.GetDB()
	if err := pool.QueryRow("getBalance", id).Scan(&bs.CustomerBalance, &bs.CreditLimit, &bs.CustomerEnabled); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("balanceStatus : ", bs)

	return bs, nil
}
