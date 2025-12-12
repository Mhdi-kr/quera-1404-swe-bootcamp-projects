package main

import (
	"log"
	"os"

	"example.com/authorization/internal/controller"
	"example.com/authorization/internal/repository"
	"example.com/authorization/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("cannot read .env from file system")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) == 0 {
		log.Fatal("jwt secret cannot be empty")
	}

	userRepo := repository.NewUserRepository()
	authSrv := service.NewAuthorizationService(jwtSecret, userRepo)
	userSrv := service.NewUserService(userRepo, authSrv)
	ctrl := controller.NewController(authSrv, userSrv)

	ctrl.ListenAndServe(":8080")
}
