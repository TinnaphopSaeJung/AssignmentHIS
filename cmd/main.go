package main

import (
	"his/internal/clients"
	"his/internal/config"
	"his/internal/database"
	"his/internal/handler"
	"his/internal/repository"
	"his/internal/routes"
	"his/internal/service"
	"his/pkg/utils"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPostgres(cfg)

	jwtManager := utils.NewJWTManager(cfg.JWTSecret)

	// init layers
	staffRepo := repository.NewStaffRepository(db)
	authService := service.NewAuthService(staffRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	hospitalAClient := clients.NewHospitalAClient()

	patientRepo := repository.NewPatientRepository(db)
	patientService := service.NewPatientService(patientRepo, hospitalAClient)
	patientHandler := handler.NewPatientHandler(patientService)

	handlers := &routes.Handlers{
		Auth:    authHandler,
		Patient: patientHandler,
	}

	r := routes.SetupRouter(handlers, jwtManager)

	r.Run(":" + cfg.AppPort)
}
