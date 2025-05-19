package handler

import (
	"net/http"

	"github.com/ShekleinAleksey/top-places/internal/entity"
	"github.com/gin-gonic/gin"
)

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
func (h *Handler) GetCountry(c *gin.Context) {
	country, err := h.service.CountryService.GetCountries()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
func (h *Handler) AddCountry(c *gin.Context) {
	var country entity.Country

	if err := c.BindJSON(&country); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CountryService.AddCountry(&country)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
