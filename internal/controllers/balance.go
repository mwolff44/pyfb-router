package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"github.com/mwolff44/pyfb-router/internal/models"
)

// BalanceController define a balance struct
type BalanceController struct{}

var balanceStatusModel = new(models.BalanceStatus)

// CheckBalance verifies is the customer has enough money to make a call
func (b BalanceController) CheckBalance(c *gin.Context) {
	// To test this, do in cli :
	// curl -i -X GET "http://127.0.0.1:8001/v1/checkbalance?customer_id=1"

	// Init value
	status := true

	// get cutomerId
	customerID, _ := c.GetQuery("customer_id")

	// Check is Query values are not empty
	if customerID == "" {
		c.JSON(400, gin.H{
			"reason": "Not all 1 expected parameter, customer_id, was present",
		})
		return
	}
	// Get balance info
	balanceStatus, err := balanceStatusModel.BalanceByCustomerID(customerID)

	switch err {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(err)
		// if customer not found for CustomerID, end of call
		c.JSON(404, gin.H{
			"reason": "No customer balance info was found",
		})
		return
	default:
		c.JSON(500, gin.H{
			"reason": "Internal server error",
		})
		return
	case nil:
		log.Println("balanceStatus : ", balanceStatus)
	}

	if balanceStatus.CustomerBalance < balanceStatus.CreditLimit {
		status = false
	}

	// Send answer
	c.JSON(200, gin.H{
		"statusBalance":  status,
		"customerStatus": balanceStatus.CustomerEnabled,
		"balance":        balanceStatus.CustomerBalance,
		"credit_limit":   balanceStatus.CreditLimit,
		"version":        "1.0",
	})

}
