package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"example.com/authorization/internal/controller"
	"example.com/authorization/internal/grpcserver"
	"example.com/authorization/internal/repository"
	"example.com/authorization/internal/service"
	"example.com/authorization/pkg"
	_ "github.com/go-sql-driver/mysql"
)

const defaultListenAddr = "0.0.0.0:3030"

func main() {
	cfg, err := pkg.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	initLogger(cfg.LogLevel)

	ctrl, postSrv, err := buildController(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := grpcserver.ListenAndServe(cfg.GrpcAddr, postSrv); err != nil {
			log.Fatal(err)
		}
	}()

	ctrl.ListenAndServe(defaultListenAddr)
}

func initLogger(level slog.Level) {
	log.SetFlags(log.Lshortfile)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})))
}

func buildController(cfg pkg.Config) (controller.Controller, service.PostService, error) {
	sqldb, err := pkg.NewSQLRepository(cfg.DBConnectionURI)
	if err != nil {
		return controller.Controller{}, service.PostService{}, fmt.Errorf("database connection failed: %w", err)
	}

	cache := pkg.NewCache(cfg.RedisAddr)

	commentRepo := repository.NewCommentRepo(sqldb, cache)
	postRepo := repository.NewPostRepository(sqldb, cache)
	userRepo := repository.NewUserRepository(sqldb)

	analyticsSrv := service.NewAnalyticsService(cache)
	authSrv := service.NewAuthorizationService(cfg.JwtSecret, userRepo)
	userSrv := service.NewUserService(userRepo, authSrv)
	postSrv := service.NewPostService(postRepo)
	commentSrv := service.NewCommentService(commentRepo)

	ctrl := controller.NewController(cfg, authSrv, userSrv, postSrv, commentSrv, analyticsSrv)

	return ctrl, postSrv, nil
}
