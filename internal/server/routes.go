package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mwolff44/pyfb-router/internal/controllers"
)

// NewRouter initializes a new gin router
func NewRouter() *gin.Engine {

	// Set the router as the default one provided by Gin
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	//router.Use(middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	{

		// The request responds to a url matching:  /v1/outboundroute?r_uri=33679590000&f_uri=33679590000&customer_id=1&socket=127.0.0.1:5060
		outboundroute := new(controllers.OutboundRouteController)
		v1.GET("/outboundroute", outboundroute.AvailableRoutes)

		// The request responds to a url matching:  /v1/balance?customer_id=1
		balance := new(controllers.BalanceController)
		v1.GET("/balance", balance.CheckBalance)

		// The request responds to a url matching:  /v1/direction?phonenumber=33240760000
		direction := new(controllers.DirectionController)
		v1.GET("/direction", direction.CheckDirection)

		// The request responds to a url matching:  /v1/customerrate?f_uri=33240760000&r_uri=33679590000&customer_id=1
		customerrate := new(controllers.CustomerRateController)
		v1.GET("/customerrate", customerrate.CheckCustomerRate)

		// The request responds to a url matching:  /v1/ratecard?customer_id=1&caller_destination_id=3
		ratecard := new(controllers.RatecardController)
		v1.GET("/ratecard", ratecard.CheckRatecard)

	}

	return router
}
