package handler

import (
	"net/http"

	"his/internal/dto"
	"his/internal/service"
	"his/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) CreateStaff(c *gin.Context) {
	var req dto.CreateStaffRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("Failed. Cannot create new staff. Error: "+err.Error(), 400))
		return
	}

	input := dto.CreateStaffInput{
		Username:   req.Username,
		Password:   req.Password,
		HospitalID: req.HospitalID,
	}

	err := h.service.CreateStaff(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error("Failed. Cannot create new staff. Error: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, utils.Success("Create new staff successfully.", nil))
}
