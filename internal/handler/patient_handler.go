package handler

import (
	"context"
	"net/http"

	"his/internal/dto"
	"his/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PatientService interface {
	Search(ctx context.Context, hospitalID int64, req dto.SearchPatientRequest) (*dto.SearchPatientResponse, int, error)
	SearchFromHISExternal(ctx context.Context, id string) (*dto.HospitalAPatientResponse, int, error)
}

type PatientHandler struct {
	service PatientService
}

func NewPatientHandler(service PatientService) *PatientHandler {
	return &PatientHandler{
		service: service,
	}
}

func (h *PatientHandler) Search(c *gin.Context) {
	var req dto.SearchPatientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("Invalid request. Error: "+err.Error()))
		return
	}

	hospitalID := c.GetInt64("hospital_id")

	results, statusCode, err := h.service.Search(c.Request.Context(), hospitalID, req)
	if err != nil {
		c.JSON(statusCode, utils.Error("Cannot search patient. Error: "+err.Error()))
		return
	}

	c.JSON(statusCode, utils.Success("Search patient successfully.", results))
}

func (h *PatientHandler) SearchFromHISExternal(c *gin.Context) {
	var req dto.SearchPatientFromHISRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("Invalid request. Error: "+err.Error()))
		return
	}

	res, statusCode, err := h.service.SearchFromHISExternal(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(statusCode, utils.Error("Cannot search patient from External HIS. Error: "+err.Error()))
		return
	}

	c.JSON(statusCode, utils.Success("Search patient from External HIS successfully.", res))
}
