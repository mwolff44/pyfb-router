package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"github.com/mwolff44/pyfb-router/internal/models"
)

// DirectionController define a balance struct
type DirectionController struct{}

var directionModel = new(models.Direction)

// CheckDirection determines the direction for a E.164 phone number (without +)
func (d DirectionController) CheckDirection(c *gin.Context) {
	// To test this, do in cli :
	// curl -i -X GET "http://127.0.0.1:8001/v1/direction?phonenumber=33240760000"

	// get numbers
	phonenumber, _ := c.GetQuery("phonenumber")

	// Check is Query values are not empty
	if phonenumber == "" {
		c.JSON(400, gin.H{
			"reason": "Not 1 expected parameters, phonenumber, was present",
		})
		return
	}

	// Get destination for Callee
	direction, err := directionModel.DirectionByPhoneNumberID(phonenumber)

	switch err {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(err)
		// if direction not found for the phonenumber
		c.JSON(404, gin.H{
			"reason": "No destination was found",
		})
		return
	default:
		c.JSON(500, gin.H{
			"reason": "Internal server error",
		})
		return
	case nil:
		log.Println("destination name : ", direction)
	}

	// Send answer
	c.JSON(200, gin.H{
		"prefix_id":        direction.PrefixID,
		"destination_id":   direction.DestinationID,
		"destination_name": direction.DestinationName,
		"version":          "1.0",
	})

}
