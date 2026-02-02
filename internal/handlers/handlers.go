package handlers

import (
	"net/http"

	"errors"

	"github.com/SEIEKSHION/Exchanger/internal/models"
	"github.com/gin-gonic/gin"
)

var (
	ValuteNotFound               = errors.New("Валюта не найдена")
	FailGettingtNumericVunitRate = errors.New("Числовое значение не было получено")
)

type Handler struct {
	valutes []models.Valute
}

func NewHandler(valutes []models.Valute) *Handler {
	return &Handler{valutes: valutes}
}

func (h *Handler) GetValutes(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, h.valutes)
}

func (h *Handler) MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

type ConvertRequest struct { // формат запроса по валюте
	Currency string  `json:"currency"`
	Quantity float64 `json:"quantity"`
}

type ConvertResponse struct { // формат овтета на запрос о валюте
	QuantityInRubles float64 `json:"quantity"`
	Error            string  `json:"error,omitempty"`
}

func (h *Handler) ConvertCurrency(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.performConversion(req.Currency, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ConvertResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{
		QuantityInRubles: result,
	})
}

func (h *Handler) performConversion(currency string, quantity float64) (float64, error) {
	valute, err := models.GetValuteByName(h.valutes, currency)
	if err != nil {
		return 0.0, ValuteNotFound
	}
	var numericVunitRate float64
	numericVunitRate, err = valute.GetNumericVunitRate()

	if err != nil {
		return 0.0, FailGettingtNumericVunitRate
	}
	return quantity * numericVunitRate, nil
}
