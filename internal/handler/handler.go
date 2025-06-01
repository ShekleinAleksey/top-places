package handler

import (
	"net/http"

	"github.com/ShekleinAleksey/top-places/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	countryHandler *CountryHandler
	placeHandler   *PlaceHandler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		countryHandler: NewCountryHandler(services.CountryService),
		placeHandler:   NewPlaceHandler(services.PlaceService),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	country := router.Group("/countries")
	{
		country.GET("/", h.countryHandler.GetCountry)
		country.GET("/:id", h.countryHandler.GetCountryByID)
		country.POST("/", h.countryHandler.AddCountry)
		country.PUT("/:id", h.countryHandler.UpdateCountry)
		country.DELETE("/:id", h.countryHandler.DeleteCountry)

		country.GET("/search", h.countryHandler.SearchCountries)

		country.GET("/:id/places", h.placeHandler.GetPlacesByCountryHandler)
	}
	places := router.Group("/places")
	{
		places.POST("/", h.placeHandler.CreatePlace)
		places.GET("/", h.placeHandler.GetAllPlaces)
		places.GET("/:id", h.placeHandler.GetPlace)
		places.PUT("/:id", h.placeHandler.UpdatePlace)
		places.DELETE("/:id", h.placeHandler.DeletePlace)

		places.GET("/search", h.placeHandler.SearchPlaces)
	}

	return router
}
