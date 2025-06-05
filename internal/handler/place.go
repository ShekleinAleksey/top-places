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

// CreatePlace создает новое место
// @Summary Создать новое место
// @Description Добавляет новое место
// @Tags Places
// @Accept json
// @Produce json
// @Param place body entity.Place true "Данные места"
// @Success 201 {object} entity.Place "Созданное место"
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /places/ [post]
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

// GetPlace возвращает место по ID
// @Summary Получить место по ID
// @Description Возвращает место по ID
// @Tags Places
// @Accept json
// @Produce json
// @Param id path int true "ID места"
// @Success 200 {object} entity.Place "Запрошенное место"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Место не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /places/{id} [get]
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

// GetAllPlaces возвращает все места
// @Summary Получить все места
// @Description Возвращает список всех мест
// @Tags Places
// @Accept json
// @Produce json
// @Success 200 {array} entity.Place "Список мест"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /places/ [get]
func (h *PlaceHandler) GetAllPlaces(c *gin.Context) {
	places, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}

// UpdatePlace обновляет данные места
// @Summary Обновить место
// @Description Обновляет информацию о месте по его ID
// @Tags Places
// @Accept json
// @Produce json
// @Param id path int true "ID места"
// @Param place body entity.Place true "Обновленные данные места"
// @Success 200 {object} entity.Place "Обновленное место"
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Failure 404 {object} map[string]string "Место не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /places/{id} [put]
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

// DeletePlace удаляет место
// @Summary Удалить место
// @Description Удаляет место по его ID
// @Tags Places
// @Accept json
// @Produce json
// @Param id path int true "ID места"
// @Success 204 "Место успешно удалено"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 404 {object} map[string]string "Место не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /places/{id} [delete]
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

// SearchPlaces ищет места по запросу
// @Summary Поиск мест
// @Description Поиск мест по названию
// @Tags Places
// @Accept json
// @Produce json
// @Param q query string false "Поисковый запрос"
// @Param limit query int false "Лимит результатов (по умолчанию 10)"
// @Success 200 {array} entity.Place "Список найденных мест"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /places/search [get]
func (h *PlaceHandler) SearchPlaces(c *gin.Context) {
	query := c.Query("q")

	limit := 10
	if l := c.Query("limit"); l != "" {
		if l, err := strconv.Atoi(l); err == nil && l > 0 {
			limit = l
		}
	}

	places, err := h.service.SearchPlaces(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, places)
}
