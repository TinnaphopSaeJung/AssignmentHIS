package handler

import (
	"context"
	"net/http"

	"his/internal/dto"
	"his/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	CreateStaff(ctx context.Context, input dto.CreateStaffInput) (int, error)
	Login(ctx context.Context, username, password, ipAddress string) (*dto.LoginResponse, int, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, int, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
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

	clientIP := c.ClientIP()

	res, statusCode, err := h.service.Login(c.Request.Context(), req.Username, req.Password, clientIP)
	if err != nil {
		c.JSON(statusCode, utils.Error("Login failed. Error: "+err.Error()))
		return
	}

	c.JSON(statusCode, utils.Success("Login successfully.", res))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("Failed. Invalid request. Error: "+err.Error()))
		return
	}

	res, statusCode, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(statusCode, utils.Error("Refresh token failed. Error: "+err.Error()))
		return
	}

	c.JSON(statusCode, utils.Success("Refresh token successfully.", res))
}
