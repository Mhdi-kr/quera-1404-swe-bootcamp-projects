package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

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

	redisAddress := os.Getenv("REDIS_ADDR")
	if len(jwtSecret) == 0 {
		log.Fatal("redis address cannot be empty")
	}

	ll := os.Getenv("LOG_LEVEL")
	logLevel := -4
	if len(ll) > 0 {
		parsedLl, _ := strconv.ParseInt(ll, 10, 64)
		logLevel = int(parsedLl)
	}

	log.SetFlags(log.Lshortfile)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(logLevel),
	})))

	slog.Debug("main", "debug", "1")
	slog.Info("main", "info", "1")
	slog.Error("main", "error", "1")

	sqldb, err := pkg.NewSQLRepository(dbConnectionUri)
	if err != nil {
		log.Fatal("database conenction failed")
	}

	cache := pkg.NewCache(redisAddress)

	commentRepo := repository.NewCommentRepo(sqldb, cache)
	postRepo := repository.NewPostRepository(sqldb, cache)
	userRepo := repository.NewUserRepository(sqldb)

	analyticsSrv := service.NewAnalyticsService(cache)
	authSrv := service.NewAuthorizationService(jwtSecret, userRepo)
	userSrv := service.NewUserService(userRepo, authSrv)
	postSrv := service.NewPostService(postRepo)
	commentSrv := service.NewCommentService(commentRepo)

	ctrl := controller.NewController(authSrv, userSrv, postSrv, commentSrv, analyticsSrv)

	ctrl.ListenAndServe("0.0.0.0:8080")
}
