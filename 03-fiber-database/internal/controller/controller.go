package controller

import (
	"context"
	"strconv"

	"example.com/authorization/internal/constants"
	"example.com/authorization/internal/service"
	"github.com/gofiber/fiber/v3"
)

type Controller struct {
	app     *fiber.App
	authSrv service.AuthService
	userSrv service.UserService
	postSrv service.PostService
}

func (ctrl Controller) ListenAndServe(addr string) {
	ctrl.app.Listen(addr)
}

func NewController(authSrv service.AuthService, userSrv service.UserService, postSrv service.PostService) Controller {
	app := fiber.New()

	ctrl := Controller{
		app:     app,
		authSrv: authSrv,
		userSrv: userSrv,
		postSrv: postSrv,
	}

	api := ctrl.app.Group("/api", func(c fiber.Ctx) error {
		return c.Next()
	})

	v1 := api.Group("/v1", func(c fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	v1profileAuthorized := v1.Group("/profile", func(c fiber.Ctx) error {
		headers := c.GetReqHeaders()
		authTokens, ok := headers["Authorization"]
		if len(authTokens) == 0 {
			return c.SendStatus(fiber.StatusForbidden)
		}

		jwtToken := authTokens[0]
		if !ok || len(jwtToken) == 0 {
			return c.SendStatus(fiber.StatusForbidden)
		}

		token, err := authSrv.ValidateToken(jwtToken)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		userIDStr, err := token.Claims.GetSubject()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		contextWithUserID := context.WithValue(c.Context(), constants.UsrIDContextKey, userID)

		c.SetContext(contextWithUserID)

		return c.Next()
	})

	ctrl.app.Get("/", ctrl.HandleHello)
	v1.Post("/register", ctrl.HandleRegister)
	v1.Post("/login", ctrl.HandleLogin)
	v1.Get("/users/", ctrl.HandleListUsers)

	v1posts := v1.Group("/posts")

	v1posts.Get("/", ctrl.HandleGetAllPosts)

	v1profileAuthorized.Get("/self", ctrl.HandleSelf)

	v1profileAuthorized.Get("/posts", ctrl.HandleListProfilePosts)
	v1profileAuthorized.Post("/posts", ctrl.HandleCreateProfilePost)

	return ctrl
}
