package models

import (
	"log"

	"github.com/mwolff44/pyfb-router/internal/db"
)

// Direction model
type Direction struct {
	DestinationID   int64
	DestinationName string
	PrefixID        int64
}

// DirectionByPhoneNumberID gets direction information from DB
func (d Direction) DirectionByPhoneNumberID(phonenumber string) (*Direction, error) {
	dn := &Direction{}
	pool := db.GetDB()
	if err := pool.QueryRow("getDirection", phonenumber).Scan(&dn.PrefixID, &dn.DestinationID, &dn.DestinationName); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("direction : ", dn)

	return dn, nil
}
