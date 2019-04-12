package models

import (
	"fmt"
	"log"

	"github.com/mwolff44/pyfb-router/internal/db"
)

// CustomerRatecard model
type CustomerRatecard struct {
	TechPrefix     string
	Priority       int
	Discount       float64
	AllowNegMargin bool
	RatecardID     int64
}

// CustomerRatecardByCustomerID gets balance information from DB form customerID and CallerDestinationID
func (c CustomerRatecard) CustomerRatecardByCustomerID(id int64, callerid int64) ([]CustomerRatecard, error) {
	pool := db.GetDB()
	rows, err := pool.Query("getCustomerRateCard", id, callerid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Iterate through the result set
	crs := []CustomerRatecard{}
	var cr CustomerRatecard
	for rows.Next() {
		err = rows.Scan(&cr.TechPrefix, &cr.Priority, &cr.Discount, &cr.AllowNegMargin, &cr.RatecardID)
		if err != nil {
			fmt.Println("customer ratecard : ", cr)
		} else {
			fmt.Println(err)
		}
		crs = append(crs, cr)

	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		fmt.Println(err)
	}

	fmt.Println("Article Instance := ", crs)

	return crs, nil
}
