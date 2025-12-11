package controller

import (
	"example.com/authorization/internal/service"
	"github.com/gofiber/fiber/v3"
)

type Controller struct {
	app     *fiber.App
	authSrv service.AuthService
}

func (ctrl Controller) ListenAndServe(addr string) {
	ctrl.app.Listen(addr)
}

func NewController(authSrv service.AuthService) Controller {
	app := fiber.New()

	ctrl := Controller{
		app:     app,
		authSrv: authSrv,
	}

	ctrl.app.Get("/", ctrl.HandleHello)
	ctrl.app.Get("/api/v1/self", ctrl.HandleSelf)
	ctrl.app.Post("/api/v1/register", ctrl.HandleRegister)
	ctrl.app.Post("/api/v1/login", ctrl.HandleLogin)

	return ctrl
}
