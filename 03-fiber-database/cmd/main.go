package main

import (
	"log"
	"os"

	"example.com/authorization/internal/controller"
	"example.com/authorization/internal/repository"
	"example.com/authorization/internal/service"
	"example.com/authorization/pkg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("cannot read .env from file system")
	}

	dbConnectionUri := os.Getenv("MYSQL_CONNECTION_URI")
	if len(dbConnectionUri) == 0 {
		log.Fatal("connectionUri cannot be empty")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) == 0 {
		log.Fatal("jwt secret cannot be empty")
	}

	sqldb, err := pkg.NewSQLRepository(dbConnectionUri)
	if err != nil {
		log.Fatal("database conenction failed")
	}
	userRepo := repository.NewUserRepository(sqldb)
	authSrv := service.NewAuthorizationService(jwtSecret, userRepo)
	userSrv := service.NewUserService(userRepo, authSrv)
	ctrl := controller.NewController(authSrv, userSrv)

	ctrl.ListenAndServe(":8080")
}
