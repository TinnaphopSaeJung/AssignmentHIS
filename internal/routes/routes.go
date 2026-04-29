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

	// auth
	r.POST("/staff/create", h.Auth.CreateStaff)
	r.POST("/staff/login", h.Auth.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtManager))

	auth.POST("/patient/search", h.Patient.Search)

	return r
}
