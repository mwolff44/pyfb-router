package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"github.com/mwolff44/pyfb-router/internal/models"
)

// CustomerRateController define a balance struct
type CustomerRateController struct{}

var rateModel = new(models.Rate)

// CheckCustomerRate verifies is the customer has enough money to make a call
func (b CustomerRateController) CheckCustomerRate(c *gin.Context) {
	// To test this, do in cli :
	// curl -i -X GET "http://127.0.0.1:8001/v1/customerrate?f_uri=33240760000&r_uri=33679590000&customer_id=1"

	// get infos
	customerID, _ := c.GetQuery("customer_id")
	callerNumber, _ := c.GetQuery("f_uri")
	calleeNumber, _ := c.GetQuery("r_uri")

	// Check is Query values are not empty
	if customerID == "" || callerNumber == "" || calleeNumber == "" {
		c.JSON(400, gin.H{
			"reason": "Not all 3 expected parameter, f_uri, r_uri, customer_id, were present",
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

	// Get callerNumber direction
	callerDirection, errCaller := directionModel.DirectionByPhoneNumberID(callerNumber)

	switch errCaller {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(errCaller)
	default:
		c.JSON(500, gin.H{
			"reason": "Internal server error",
		})
		return
	case nil:
		log.Println("caller destination : ", callerDirection)
	}

	// get calleeNumber direction
	calleeDirection, errCallee := directionModel.DirectionByPhoneNumberID(calleeNumber)

	switch errCallee {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(errCallee)
		// if direction not found for the phonenumber
		c.JSON(404, gin.H{
			"reason": "No destination was found",
		})
	default:
		c.JSON(500, gin.H{
			"reason": "Internal server error",
		})
		return
	case nil:
		log.Println("callee destination : ", calleeDirection)
	}

	// get ratecard
	ratecard, errRC := ratecardModel.CustomerRatecardByCustomerID(custID, callerDirection.DestinationID)

	switch errRC {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(errRC)
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

	// Get rate info
	l := len(ratecard)
	if l == 0 {
		log.Println("No ratecard : END")
		// if customer ratecard not found for CustomerID and callerDestinationID, end of call
		c.JSON(404, gin.H{
			"reason": "No customer ratecard info was found",
		})
		return
	}

	log.Println("nb ratecard : ", l)
	for k, s := range ratecard {
		rate := new(models.Rate)
		rate, errRate := rateModel.CustomerRateByRatecardID(s.RatecardID, calleeNumber, calleeDirection.DestinationID)

		switch errRate {
		case pgx.ErrNoRows:
			//c.AbortWithStatus(404)
			fmt.Println(errRate)
			fmt.Println(s)
			// if rate not found for CustomerID, continue else no more ratecard
			if k == (l - 1) {
				c.JSON(404, gin.H{
					"reason": "No customer rate info was found",
				})
			}
			return
		default:
			c.JSON(500, gin.H{
				"reason": "Internal server error",
			})
			return
		case nil:
			log.Println("controllers:rate: customer rate : ", rate)
			// rate found, end of boucle
			c.JSON(200, gin.H{
				"rate":    rate,
				"version": "1.0",
			})
			return
		}
	}
}
