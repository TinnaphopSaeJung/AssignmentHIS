package routes

import (
	"his/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *handler.AuthHandler
	// Patient *handler.PatientHandler
}

func SetupRouter(h *Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// auth
	r.POST("/staff/create", h.Auth.CreateStaff)
	r.POST("/staff/login", h.Auth.Login)

	// patient
	// r.GET("/patient/search", h.Patient.Search)

	return r
}
