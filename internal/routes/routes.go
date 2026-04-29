package routes

import (
	"his/internal/handler"
	"his/internal/middleware"
	"his/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth    *handler.AuthHandler
	Patient *handler.PatientHandler
}

func SetupRouter(h *Handlers, jwtManager *utils.JWTManager) *gin.Engine {
	r := gin.Default()

	// staff
	staff := r.Group("/staff")
	{
		staff.POST("/create", h.Auth.CreateStaff)
		staff.POST("/login", h.Auth.Login)
	}

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtManager))

	// patient
	patient := auth.Group("/patient")
	{
		patient.POST("/search", h.Patient.Search)
	}

	return r
}
