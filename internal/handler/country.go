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

// @Summary GetCountry
// @Tags country
// @Description get country
// @ID get-country
// @Accept  json
// @Produce  json
// @Param id path string true "Country ID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /country/ [get]
func (h *CountryHandler) GetCountry(c *gin.Context) {
	country, err := h.service.GetCountries()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, country)
}

// @Summary Get country by ID
// @Tags country
// @Description Get country by ID
// @ID get-country-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "Country ID"
// @Success 200 {object} entity.Country
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /country/{id} [get]
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
// @Tags country
// @Description add country
// @ID add-country
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.User
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /country/ [post]
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

// @Summary Delete country
// @Tags country
// @Description delete country by ID
// @ID delete-country
// @Accept  json
// @Produce  json
// @Param id path int true "Country ID"
// @Success 200 {object} map[string]interface{} "{"status": "success", "deleted_id": id}"
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /country/{id} [delete]
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
