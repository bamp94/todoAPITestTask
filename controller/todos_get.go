package controller

import (
	"net/http"

	"cyberzilla_api_task/model"

	"github.com/labstack/echo/v4"
)

// @Summary Проверка сервера
// @Description Запускает процесс проверки всех зависимостей, необходимых для корректной работы сервера
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @tags Служебные
// @in header
// @Param token query string false "Токен списка дел" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Success 200 {object} HealthCheck
// @Router /todos [get]
func (c *Controller) getTodoList(ctx echo.Context) error {
	auth := ctx.Request().Header.Get("Authorization")
	token := ctx.QueryParam("token")
	var (
		res []model.TodoTask
		err error
	)
	switch {
	case auth != "":
		res, err = c.app.TodoListTasks(auth)
	case token != "":
		res, err = c.app.TodoListTasks(token)
	default:
		res, err = c.app.TodoListTasks("")
	}
	if err != nil {
		return c.respondError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, echo.Map{"result": res})
}
