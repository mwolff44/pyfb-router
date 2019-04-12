package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"github.com/mwolff44/pyfb-router/internal/models"
)

// RatecardController define a balance struct
type RatecardController struct{}

var ratecardModel = new(models.CustomerRatecard)

// CheckRatecard verifies is the customer has enough money to make a call
func (r RatecardController) CheckRatecard(c *gin.Context) {
	// To test this, do in cli :
	// curl -i -X GET "http://127.0.0.1:8001/v1/checkbalance?customer_id=1&caller_destination_id=3"

	// get cutomerId and caller_destnation_id
	customerID, _ := c.GetQuery("customer_id")
	callerDestinationID, _ := c.GetQuery("caller_destination_id")

	// Check is Query values are not empty
	if customerID == "" || callerDestinationID == "" {
		c.JSON(400, gin.H{
			"reason": "Not all 2 expected parameter, customer_id, caller_destination_id, were present",
		})
		return
	}
	custID, errCID := strconv.ParseInt(customerID, 10, 64)
	if errCID != nil {
		c.JSON(400, gin.H{
			"reason": "customer_id param is not an integer",
		})
		return
	}
	cdID, errCID := strconv.ParseInt(callerDestinationID, 10, 64)
	if errCID != nil {
		c.JSON(400, gin.H{
			"reason": "caller_destination_id param is not an integer",
		})
		return
	}

	// Get balance info
	ratecard, err := ratecardModel.CustomerRatecardByCustomerID(custID, cdID)

	switch err {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(err)
		// if customer ratecard not found for CustomerID and callerDestinationID, end of call
		c.JSON(404, gin.H{
			"reason": "No customer ratecard info was found",
		})
		return
	default:
		c.JSON(500, gin.H{
			"reason": "Internal server error",
		})
		return
	case nil:
		log.Println("ratecard : ", ratecard)
	}

	// Send answer
	c.JSON(200, gin.H{
		"ratecard": ratecard,
		"version":  "1.0",
	})

}
