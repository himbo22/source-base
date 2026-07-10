package bootstrap

import (
	"source-base/internal/ports"
	"source-base/pkg/common/http/handler"
	"source-base/pkg/common/http/middlewares"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

type Controllers struct {
	UserController ports.UserController
}

func InitRouter(
	userController ports.UserController,
) *Controllers {
	return &Controllers{
		UserController: userController,
	}
}

func (c *Controllers) RegisterRouters(echo *echo.Echo) {
	api := echo.Group("/api/v1/public")

	// user api
	users := api.Group("/users")
	{
		users.POST("", handler.Wrap(c.UserController.Create))
		users.GET("/:public_id", handler.Wrap(c.UserController.GetByPublicID))
	}
}

func InitEcho(controllers *Controllers, logger *zap.Logger) *echo.Echo {
	e := echo.New()

	e.Use(echoMiddleware.RequestID())
	e.Use(middlewares.CustomRecover(logger))
	e.Use(middlewares.RequestLogger(logger))
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	//e.IPExtractor = echo.ExtractIPFromXFFHeader()

	errorHandler := handler.NewErrorHandler(logger)

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		errorHandler.ErrorHandler(c, err)
	}

	e.GET("/ping", Ping)

	// register all controllers
	controllers.RegisterRouters(e)

	return e
}

func Ping(c *echo.Context) error {
	err := c.JSON(200, "pong")
	if err != nil {
		return err
	}
	return nil
}
