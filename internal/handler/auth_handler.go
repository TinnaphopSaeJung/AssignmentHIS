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
		c.JSON(http.StatusBadRequest, utils.Error("Failed. Cannot create new staff. Error: "+err.Error()))
		return
	}

	input := dto.CreateStaffInput{
		Username:   req.Username,
		Password:   req.Password,
		HospitalID: req.HospitalID,
	}

	statusCode, err := h.service.CreateStaff(c.Request.Context(), input)
	if err != nil {
		c.JSON(statusCode, utils.Error(err.Error()))
		return
	}

	c.JSON(statusCode, utils.Success("Create new staff successfully.", nil))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("Failed. Invalid request. Error: "+err.Error()))
		return
	}

	res, statusCode, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(statusCode, utils.Error("Login failed. Error: "+err.Error()))
		return
	}

	c.JSON(statusCode, utils.Success("Login successfully.", res))
}
