package models

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/pgtype"

	"github.com/mwolff44/pyfb-router/internal/db"
)

// Route handle routing informations
type Route struct {
	ID                 int64
	DestnumLengthMap   int
	RouteType          int
	ProviderEndpointID int64
	ProviderRatecardID int64
	RouteRule          string
	Status             string
	Weight             int
	Priority           int
	Prefix             pgtype.Varchar
	DestnumLength      pgtype.Int4
	DestinationID      pgtype.Int4
	CountryID          pgtype.Int4
	TypeID             pgtype.Int4
	RegionID           pgtype.Int4
	Name               string
	CalleeNormID       pgtype.Int4
	CalleeNormInID     pgtype.Int4
	CalleridInFrom     bool
	CalleridNormID     pgtype.Int4
	CalleridNormInID   pgtype.Int4
	FromDomain         string
	Pai                bool
	Ppi                bool
	Pid                bool
	GwPrefix           string
	SipTransport       string
	SipPort            int
	SipProxy           string
	Username           string
	Suffix             string
	EstimatedQuality   int
	RatecardPrefix     string
	ProviderRate       int64
}

// RouteDetail contains all data of one route
type RouteDetail struct {
	BranchFlags int    `json:"branch_flags"`
	DstURI      string `json:"dst_uri"`
	FrInvTimer  int    `json:"fr_inv_timer"`
	FrTimer     int    `json:"fr_timer"`
	Headers     Header `json:"headers"`
	Path        string `json:"path,omitempty"`
	Socket      string `json:"socket"`
	URI         string `json:"uri"`
}

// Header contains from and to headers data
type Header struct {
	Extra string       `json:"extra"`
	From  HeaderDetail `json:"from"`
	To    HeaderDetail `json:"to"`
}

// HeaderDetail contains Display and uri fields
type HeaderDetail struct {
	Display string `json:"display"`
	URI     string `json:"uri"`
}

// OutBoundRouteByCustomerID gets routes information from DB for customerID
func (o Route) OutBoundRouteByCustomerID(customerID int64, calleeID string, calleeDestinationID int64, callerID string, socket string) ([]Route, error) {
	pool := db.GetDB()
	rows, err := pool.Query("getOutboundRoutes", customerID, calleeID, calleeDestinationID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var custRoutes []Route

	// Iterate through the result set
	var r Route
	for rows.Next() {
		err = rows.Scan(
			&r.ID,
			&r.DestnumLengthMap,
			&r.RouteType,
			&r.ProviderEndpointID,
			&r.ProviderRatecardID,
			&r.RouteRule,
			&r.Status,
			&r.Weight,
			&r.Priority,
			&r.Prefix,
			&r.DestnumLength,
			&r.DestinationID,
			&r.CountryID,
			&r.TypeID,
			&r.RegionID,
			&r.Name,
			&r.CalleeNormID,
			&r.CalleeNormInID,
			&r.CalleridInFrom,
			&r.CalleridNormID,
			&r.CalleridNormInID,
			&r.FromDomain,
			&r.Pai,
			&r.Ppi,
			&r.Pid,
			&r.GwPrefix,
			&r.SipTransport,
			&r.SipPort,
			&r.SipProxy,
			&r.Username,
			&r.Suffix,
			&r.EstimatedQuality,
			&r.RatecardPrefix)
		fmt.Println("route : ", r)
		if err != nil {
			fmt.Println("route : ", r)
		} else {
			fmt.Println(err)
		}
		custRoutes = append(custRoutes, r)

	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return nil, rows.Err()
	}
	fmt.Println("route list : ", custRoutes)

	return custRoutes, nil
}
