package main

import (
	"his/internal/config"
	"his/internal/database"
	"his/internal/handler"
	"his/internal/repository"
	"his/internal/routes"
	"his/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPostgres(cfg)

	// init layers
	staffRepo := repository.NewStaffRepository(db)
	authService := service.NewAuthService(staffRepo)
	authHandler := handler.NewAuthHandler(authService)

	handlers := &routes.Handlers{
		Auth: authHandler,
		// Patient: patientHandler, // เดี๋ยวเราจะทำ
	}

	r := routes.SetupRouter(handlers)

	r.Run(":" + cfg.AppPort)
}
