package main

import (
	"fmt"
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

const defaultListenAddr = "0.0.0.0:3030"

type config struct {
	dbConnectionURI string
	jwtSecret       string
	redisAddress    string
	logLevel        slog.Level
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	initLogger(cfg.logLevel)

	ctrl, err := buildController(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctrl.ListenAndServe(defaultListenAddr)
}

func loadConfig() (config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return config{}, fmt.Errorf("cannot read .env from file system: %w", err)
	}

	dbConnectionURI, err := requireEnv("MYSQL_CONNECTION_URI")
	if err != nil {
		return config{}, err
	}

	jwtSecret, err := requireEnv("JWT_SECRET")
	if err != nil {
		return config{}, err
	}

	redisAddress, err := requireEnv("REDIS_ADDR")
	if err != nil {
		return config{}, err
	}

	logLevel := slog.LevelDebug
	if ll := os.Getenv("LOG_LEVEL"); ll != "" {
		parsedLl, err := strconv.ParseInt(ll, 10, 64)
		if err != nil {
			return config{}, fmt.Errorf("invalid LOG_LEVEL %q: %w", ll, err)
		}
		logLevel = slog.Level(parsedLl)
	}

	return config{
		dbConnectionURI: dbConnectionURI,
		jwtSecret:       jwtSecret,
		redisAddress:    redisAddress,
		logLevel:        logLevel,
	}, nil
}

func requireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s cannot be empty", key)
	}

	return value, nil
}

func initLogger(level slog.Level) {
	log.SetFlags(log.Lshortfile)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})))
}

func buildController(cfg config) (controller.Controller, error) {
	sqldb, err := pkg.NewSQLRepository(cfg.dbConnectionURI)
	if err != nil {
		return controller.Controller{}, fmt.Errorf("database connection failed: %w", err)
	}

	cache := pkg.NewCache(cfg.redisAddress)

	commentRepo := repository.NewCommentRepo(sqldb, cache)
	postRepo := repository.NewPostRepository(sqldb, cache)
	userRepo := repository.NewUserRepository(sqldb)

	analyticsSrv := service.NewAnalyticsService(cache)
	authSrv := service.NewAuthorizationService(cfg.jwtSecret, userRepo)
	userSrv := service.NewUserService(userRepo, authSrv)
	postSrv := service.NewPostService(postRepo)
	commentSrv := service.NewCommentService(commentRepo)

	ctrl := controller.NewController(authSrv, userSrv, postSrv, commentSrv, analyticsSrv)

	return ctrl, nil
}
