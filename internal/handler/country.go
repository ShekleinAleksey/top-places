package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/ShekleinAleksey/top-places/internal/service"
	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	service *service.CountryService
}

func NewCountryHandler(service *service.CountryService) *CountryHandler {
	return &CountryHandler{service: service}
}

// @Summary Get all countries
// @Tags Countries
// @Description Retrieve a list of all countries
// @ID get-country
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Country
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /countries/ [get]
func (h *CountryHandler) GetCountry(c *gin.Context) {
	country, err := h.service.GetCountries()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, country)
}

// @Summary Get country by ID
// @Tags Countries
// @Description Get country by ID
// @ID get-country-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "Country ID"
// @Success 200 {object} entity.Country
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /countries/{id} [get]
func (h *CountryHandler) GetCountryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	country, err := h.service.GetCountryByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, country)
}

// @Summary AddCountry
// @Tags Countries
// @Description add country
// @ID add-country
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Country
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /countries/ [post]
func (h *CountryHandler) AddCountry(c *gin.Context) {
	var country entity.Country

	if err := c.BindJSON(&country); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.AddCountry(&country)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// UpdateCountry godoc
// @Summary Обновить страну
// @Description Обновляет данные страны по ID
// @Tags Countries
// @Accept json
// @Produce json
// @Param id path int true "ID страны"
// @Param country body entity.Country true "Данные для обновления"
// @Success 200 {object} entity.Country
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries/{id} [put]
func (h *CountryHandler) UpdateCountry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country ID"})
		return
	}

	var country entity.Country
	if err := c.ShouldBindJSON(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	country.ID = id

	updatedCountry, err := h.service.UpdateCountry(&country)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "country not found" {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "is required") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCountry)
}

// @Summary Delete country
// @Tags Countries
// @Description delete country by ID
// @ID delete-country
// @Accept  json
// @Produce  json
// @Param id path int true "Country ID"
// @Success 200 {object} map[string]interface{} "{"status": "success", "deleted_id": id}"
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /countries/{id} [delete]
func (h *CountryHandler) DeleteCountry(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	deletedID, err := h.service.DeleteCountry(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "success",
		"deleted_id": deletedID,
	})
}

// SearchCountries godoc
// @Summary Search countries
// @Description Search countries by name with optional limit
// @Tags Countries
// @Accept  json
// @Produce  json
// @Param q query string true "Search query (minimum 2 characters)"
// @Param limit query int false "Maximum number of results (default: 10)"
// @Success 200 {array} entity.Country "List of matching countries"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /countries/search [get]
func (h *CountryHandler) SearchCountries(c *gin.Context) {
	query := c.Query("q")
	// if len(query) < 2 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Query must be at least 2 characters"})
	// 	return
	// }

	limit := 10
	if l := c.Query("limit"); l != "" {
		if l, err := strconv.Atoi(l); err == nil && l > 0 {
			limit = l
		}
	}

	countries, err := h.service.SearchCountries(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if countries == nil {
		countries = []entity.Country{}
	}

	c.JSON(http.StatusOK, countries)
}
