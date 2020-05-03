package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4/middleware"

	"cyberzilla_api_task/application"
	"cyberzilla_api_task/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var errWrongTodoTaskID = errors.New("Неправильно задан ID задачи")
var errEmptyAuthorizationToken = errors.New("Не задан авторизационный токен")

// Controller is presentation tier of 3-layer architecture
type Controller struct {
	app    application.Application
	config config.Main
	router *echo.Echo
}

// New Controller constructor
func New(config config.Main, app application.Application) Controller {
	return Controller{
		app:    app,
		config: config,
		router: echo.New(),
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ServeHTTP http server
func (c Controller) ServeHTTP(ctx context.Context, port int) {
	c.initRoutes()

	go func() {
		if err := c.router.Start(fmt.Sprint(":", port)); err != nil {
			logrus.WithError(err).Fatal("can't start serving http")
		}
	}()

	// Gracefully stopping
	<-ctx.Done()
	if err := c.router.Shutdown(ctx); err != nil {
		logrus.Error("http server shutdown error:", err)
	}
	logrus.Println("http server has stopped")
}

func (c Controller) initRoutes() {
	// cors used for success answer on OPTION request from swagger
	c.router.Use(middleware.CORS())

	// init validator
	c.router.Validator = &CustomValidator{validator: validator.New()}

	// init logger
	c.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	c.router.GET("/swagger/*any", echoSwagger.WrapHandler)
	c.router.GET("/healthcheck", c.healthcheck)

	c.router.GET("/todos", c.getTodoList)
	c.router.POST("/todos", c.createTodoTask)
	c.router.GET("/todos/:id", c.getTodoTask)
	c.router.PUT("/todos/:id", c.updateTodoTask)
}

func getAuthorizationToken(ctx echo.Context) (string, error) {
	auth := ctx.Request().Header.Get("Authorization")
	token := ctx.QueryParam("token")
	switch {
	case auth != "":
		return auth, nil
	case token != "":
		return token, nil
	default:
		return "", errEmptyAuthorizationToken
	}
}

func getIntParam(ctx echo.Context, name string) (int, error) {
	param, err := strconv.Atoi(ctx.Param(name))
	if err != nil {
		return 0, err
	}
	return param, nil
}
