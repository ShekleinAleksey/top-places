package handler

import (
	"net/http"
	"strconv"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/ShekleinAleksey/top-places/internal/service"
	"github.com/gin-gonic/gin"
)

type PlaceHandler struct {
	service *service.PlaceService
}

func NewPlaceHandler(service *service.PlaceService) *PlaceHandler {
	return &PlaceHandler{service: service}
}

func (h *PlaceHandler) CreatePlace(c *gin.Context) {
	var place entity.Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	createdPlace, err := h.service.Create(&place)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPlace)
}

func (h *PlaceHandler) GetPlace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	place, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "place not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, place)
}

func (h *PlaceHandler) GetAllPlaces(c *gin.Context) {
	places, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}

func (h *PlaceHandler) UpdatePlace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var place entity.Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	place.ID = id

	updatedPlace, err := h.service.Update(&place)
	if err != nil {
		if err.Error() == "place not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedPlace)
}

func (h *PlaceHandler) DeletePlace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "place not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPlacesByCountryHandler возвращает места по стране
// @Summary Получить места по стране
// @Description Возвращает список всех мест для указанной страны
// @Tags Places
// @Produce json
// @Param country_id path int true "ID страны"
// @Success 200 {array} entity.Place
// @Failure 404 {object} map[string]string
// @Router /countries/{country_id}/places [get]
func (h *PlaceHandler) GetPlacesByCountryHandler(c *gin.Context) {
	countryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country ID"})
		return
	}

	places, err := h.service.GetPlacesByCountry(countryID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "country not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}
