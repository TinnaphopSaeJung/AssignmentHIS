package main

import (
	"his/internal/clients"
	"his/internal/config"
	"his/internal/database"
	"his/internal/handler"
	"his/internal/messaging"
	"his/internal/repository"
	"his/internal/routes"
	"his/internal/service"
	"his/pkg/utils"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewPostgres(cfg)
	redisClient := database.NewRedis(cfg)

	rabbitConn := database.NewRabbitMQ(cfg)
	defer rabbitConn.Close()

	auditPublisher := messaging.NewRabbitMQAuditPublisher(rabbitConn)

	jwtManager := utils.NewJWTManager(cfg.JWTSecret)

	// init layers
	staffRepo := repository.NewStaffRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	loginAttemptService := service.NewLoginAttemptService(redisClient)

	authService := service.NewAuthService(
		staffRepo,
		refreshTokenRepo,
		jwtManager,
		loginAttemptService,
		auditPublisher,
	)

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
