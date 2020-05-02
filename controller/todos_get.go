package controller

import (
	"github.com/labstack/echo/v4"
)

// @Summary Взять список заданий
// @Description Возвращает список заданий
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @tags Основные
// @in header
// @Param token query string false "Токен списка дел" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Success 200 {object} HealthCheck
// @Router /todos [get]
func (c *Controller) getTodoList(ctx echo.Context) error {
	token := c.getAuthorizationToken(ctx)
	res, err := c.app.TodoListTasks(token)
	if err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, res)
}
