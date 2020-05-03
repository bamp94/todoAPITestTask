package controller

import (
	"cyberzilla_api_task/model"

	"github.com/labstack/echo/v4"
)

type GetTodoListResponse struct {
	Result []model.TodoTask `json:"result"`
}

// @Summary Список заданий
// @Description Возвращает список заданий
// @Accept  json
// @Produce  json
// @tags Основные
// @in header
// @Param token query string false "Токен списка задач" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Success 200 {object} GetTodoListResponse
// @Failure 400 {object} ErrBadRequest
// @Failure 404 {object} ErrNotFound
// @Failure 500 {object} ErrInternal
// @Router /todos [get]
func (c *Controller) getTodoList(ctx echo.Context) error {
	token, err := getAuthorizationToken(ctx)
	if err != nil {
		return c.respondError(ctx, err)
	}
	res, err := c.app.TodoTasksList(token)
	if err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, res)
}
