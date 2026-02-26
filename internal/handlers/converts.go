package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/SEIEKSHION/Exchanger/internal/domain"
	"github.com/gin-gonic/gin"
)

var (
	ValuteNotFound               = errors.New("Валюта не найдена")
	FailGettingtNumericVunitRate = errors.New("Числовое значение не было получено")
)

type Handler struct {
	valutes []domain.Valute
}

func NewHandler(valutes []domain.Valute) *Handler {
	return &Handler{valutes: valutes}
}

func (h *Handler) GetValutes(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, h.valutes)
}

func (h *Handler) MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) Exchanger(c *gin.Context) {
	c.HTML(http.StatusOK, "exchanger.html", nil)
}

type ConvertRequest struct { // формат запроса по валюте
	FromCurrency string  `json:"fromcurrency"`
	ToCurrency   string  `json:"tocurrency"`
	Quantity     float64 `json:"quantity"`
}

type ConvertResponse struct { // формат овтета на запрос о валюте
	Quantity float64 `json:"quantity"`
	Error    string  `json:"error,omitempty"`
}

func (h *Handler) ConvertCurrency(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// конвертируем из одной валюты во вторую
	result, err := h.convertOperation(req.FromCurrency, req.ToCurrency, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ConvertResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{
		Quantity: result,
	})
}

func RoundTo(val float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(val*p) / p
}

func (h *Handler) convertOperation(fromCurrency, toCurrency string, quantity float64) (float64, error) {
	// валюта из которой
	fromValute, err := domain.GetValuteByName(h.valutes, fromCurrency)
	if err != nil {
		return 0.0, ValuteNotFound
	}

	// валюта в которую
	var toValute domain.Valute

	toValute, err = domain.GetValuteByName(h.valutes, toCurrency)
	if err != nil {
		return 0.0, ValuteNotFound
	}

	// инициализация численных значений
	var fromRate, toRate float64

	// получение первого численного значения
	fromRate, err = fromValute.GetNumericVunitRate()
	if err != nil {
		return 0.0, FailGettingtNumericVunitRate
	}
	// получение второго численного значения
	toRate, err = toValute.GetNumericVunitRate()
	if err != nil {
		return 0.0, FailGettingtNumericVunitRate
	}
	fmt.Println("fromCurrency: ", fromCurrency)
	fmt.Println("toCurrency: ", toCurrency)

	// нужно перевести из начальной в конечную: поделить конечную на начальную, умножить на количество начальных
	result := RoundTo(quantity*fromRate/toRate, 3)

	return result, nil
}
