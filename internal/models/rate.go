package models

import (
	"log"

	"github.com/jackc/pgx/pgtype"
	"github.com/mwolff44/pyfb-router/internal/db"
)

// Rate handle rate
type Rate struct {
	ID               int            `json:"id"`
	DestnumLengthMap int            `json:"destnum_length_map"`
	RatecardID       int64          `json:"ratecard_id"`
	RateType         int            `json:"rate_type"`
	RatecardName     string         `json:"ratecard_name"`
	RCType           string         `json:"rc_type"`
	Status           string         `json:"status"`
	Rate             int64          `json:"rate"`
	BlockMinDuration int            `json:"block_min_duration"`
	MinimalTime      int            `json:"minimal_time"`
	InitBlock        int64          `json:"init_block"`
	Prefix           pgtype.Varchar `json:"prefix"`
	DestnumLength    pgtype.Int4    `json:"destnum_length"`
	DestinationID    pgtype.Int4    `json:"destination_id"`
	CountryID        pgtype.Int4    `json:"country_id"`
	TypeID           pgtype.Int4    `json:"type_id"`
	RegionID         pgtype.Int4    `json:"region_id"`
}

// CustomerRateByRatecardID gets rate informations from DB
func (c Rate) CustomerRateByRatecardID(ratecardID int64, calleeID string, calleeDestinationID int64) (*Rate, error) {
	r := &Rate{}
	pool := db.GetDB()
	if err := pool.QueryRow("getCustomerRate", ratecardID, calleeID, calleeDestinationID).Scan(
		&r.ID,
		&r.DestnumLengthMap,
		&r.RatecardID,
		&r.RateType,
		&r.RatecardName,
		&r.RCType,
		&r.Status,
		&r.Rate,
		&r.BlockMinDuration,
		&r.MinimalTime,
		&r.InitBlock,
		&r.Prefix,
		&r.DestnumLength,
		&r.DestinationID,
		&r.CountryID,
		&r.TypeID,
		&r.RegionID); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("models:rate: customer rate : ", r)

	return r, nil
}
