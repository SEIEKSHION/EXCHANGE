package handlers

import (
	"fmt"
	"net/http"

	"github.com/SEIEKSHION/Exchanger/internal/models"
	"github.com/gin-gonic/gin"
)

type MusclesHandler struct {
	muscles []models.Muscle
}

func NewMuscleHandler(muscles []models.Muscle) *MusclesHandler {
	return &MusclesHandler{muscles: muscles}
}

func (m *MusclesHandler) GetMuscles(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	fmt.Println(m.muscles)
	c.JSON(http.StatusOK, m.muscles)
}

func (m *MusclesHandler) MusclesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "muscles.html", nil)
}
