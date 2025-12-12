package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"example.com/authorization/internal/controller"
	"example.com/authorization/internal/repository"
	"example.com/authorization/internal/repository/entity"
	"example.com/authorization/internal/service"
	_ "github.com/go-sql-driver/mysql"
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

	// demo how we are going to use sql package and why it gets difficult to work with
	db, err := sql.Open("mysql", "root:example@tcp(localhost:3306)/quera-bootcamp")
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.QueryContext(context.Background(), "select username as Username, password as Passwrod from user;")

	fmt.Println(res)

	var users []entity.User

	for res.Next() {
		var u entity.User
		res.Scan(&u.Username, &u.Password)
		users = append(users, u)
	}

	fmt.Println(users)

	userRepo := repository.NewUserRepository()
	authSrv := service.NewAuthorizationService(jwtSecret, userRepo)
	userSrv := service.NewUserService(userRepo, authSrv)
	ctrl := controller.NewController(authSrv, userSrv)

	ctrl.ListenAndServe(":8080")
}
