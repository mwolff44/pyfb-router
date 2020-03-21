package controllers

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	"github.com/mwolff44/pyfb-router/internal/models"
)

// OutboundRouteController define a balance struct
type OutboundRouteController struct{}

var outboundRouteModel = new(models.Route)
var routesList = new(models.Route)

// AvailableRoutes verifies is the customer has enough money to make a call
func (o OutboundRouteController) AvailableRoutes(c *gin.Context) {
	// To test this, do in cli :
	// curl -i -X GET "http://127.0.0.1:8001/v1/outboundroute?r_uri=33679590000&f_uri=33679590000&customer_id=1&socket=127.0.0.1:5060"

	// get infos
	customerID, _ := c.GetQuery("customer_id")
	callerNumber, _ := c.GetQuery("f_uri")
	calleeNumber, _ := c.GetQuery("r_uri")
	socket, _ := c.GetQuery("socket")

	// Check is Query values are not empty
	if customerID == "" || callerNumber == "" || calleeNumber == "" || socket == "" {
		c.JSON(400, gin.H{
			"reason": "Not all 4 expected parameter, f_uri, r_uri, customer_id, socket, were present",
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
		return
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

	var custRate = new(models.Rate)
	// Get rate info
	l := len(ratecard)
	log.Println("nb ratecard : ", l)

	if l == 0 {
		log.Println("No ratecard : END")
		// if customer ratecard not found for CustomerID and callerDestinationID, end of call
		c.JSON(404, gin.H{
			"reason": "No customer ratecard info was found",
		})
		return
	}

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
				return
			}
		default:
			c.JSON(500, gin.H{
				"reason": "Internal server error",
			})
			return
		case nil:
			log.Println("controllers:rate: customer rate : ", rate)
			custRate = rate
		}
		if custRate != nil {
			break
		}
	}

	// get unordered routes

	custRoutes, errRoute := outboundRouteModel.OutBoundRouteByCustomerID(custID, calleeNumber, calleeDirection.DestinationID, callerNumber, socket)
	switch errRoute {
	case pgx.ErrNoRows:
		//c.AbortWithStatus(404)
		fmt.Println(errRoute)
		// if customer not found for CustomerID, end of call
		c.JSON(404, gin.H{
			"reason": "No route was found",
		})
		return
	default:
		c.JSON(500, gin.H{
			"reason": "Internal server error",
		})
		return
	case nil:
		log.Println("route : ", custRoutes)
	}

	// Get providerRates by route
	// think to check margin ToDo
	for i, s := range custRoutes {
		provRate := new(models.Rate)
		var errorRate error
		// Get Rate from provider ratecard
		provRate, errorRate = rateModel.ProviderRateByRatecardID(s.ProviderRatecardID, calleeNumber, calleeDirection.DestinationID)
		fmt.Println("provider rate :", provRate)
		fmt.Println("error provider rate :", errorRate)
		custRoutes[i].ProviderRate = provRate.Rate
	}
	fmt.Println("route list with prov rates : ", custRoutes)

	if custRoutes[0].RouteRule == "LCR" {
		fmt.Println("LCR rules")
		sort.Slice(custRoutes, func(i, j int) bool {
			return custRoutes[i].ProviderRate < custRoutes[j].ProviderRate
		})
	} else if custRoutes[0].RouteRule == "PRIO" {
		fmt.Println("PRIO rules")
		sort.SliceStable(custRoutes, func(i, j int) bool {
			return custRoutes[i].Priority < custRoutes[j].Priority
		})
		fmt.Println("route list ordered :", custRoutes)
	} else if custRoutes[0].RouteRule == "QUALITY" {
		fmt.Println("QUALITY rules")
		sort.SliceStable(custRoutes, func(i, j int) bool {
			return custRoutes[i].EstimatedQuality > custRoutes[j].EstimatedQuality
		})
		fmt.Println("route list ordered :", custRoutes)
	} else if custRoutes[0].RouteRule == "WEIGHT" {
		// Calculate total weight
		var totalWeight int
		for _, c := range custRoutes {
			totalWeight += c.Weight
		}
		//ToDo
		// example of algo : https://medium.com/@peterkellyonline/weighted-random-selection-3ff222917eb6
		// for the moment, weight is not taken, only fair distribution
	}
	fmt.Println("route list ordered :", custRoutes)

	//var routesList
	var routesList []models.RouteDetail

	// Handles the creation of route for all available routes and gateways
	for _, s := range custRoutes {
		// Handle CallerID rules

		newCallerID, headersExtra := models.SetCallerID(callerNumber, s.Pai, s.Ppi, s.Pid, s.AddPlusInCaller, socket)
		fmt.Println("newCallerID :", newCallerID)

		// Handle Callee rules
		newCalleeID := models.SetCallee(calleeNumber, s.GwPrefix, s.Suffix, s.RatecardPrefix)
		fmt.Println("newCalleeID :", newCalleeID)

		// Contruct routes structures
		toHeader := models.HeaderDetail{Display: newCalleeID, URI: "sip:" + newCalleeID + "@" + s.SipProxy}

		fromHeader := models.ConstructFromHeader(s.CalleridInFrom, newCallerID, socket, s.Username)

		headers := models.Header{Extra: headersExtra, From: fromHeader, To: toHeader}

		routeDetail := models.RouteDetail{URI: "sip:" + newCalleeID + "@" + s.SipProxy,
			DstURI:      "sip:" + s.SipProxy + ":" + strconv.Itoa(s.SipPort),
			Socket:      socket,
			Headers:     headers,
			BranchFlags: 8,
			FrTimer:     5000,
			FrInvTimer:  30000}

		routesList = append(routesList, routeDetail)
	}

	// ToDo : add 2 objects : direction for caller and callee

	c.JSON(200, gin.H{
		"rate":    custRate,
		"routes":  routesList,
		"routing": "serial",
		"version": "1.0",
	})

}
