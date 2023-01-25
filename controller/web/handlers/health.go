package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get health check
type (
	HealthGetResponse struct {
		Message string `json:"message"`
	}
	HealthGetHandler struct{}
)

func NewHealthGetHandler() *HealthGetHandler {
	return &HealthGetHandler{}
}

func (h *HealthGetHandler) Handle(gc *gin.Context) {
	response := HealthGetResponse{
		Message: "OK",
	}
	gc.JSON(http.StatusOK, response)
}
