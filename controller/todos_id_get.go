package controller

import (
	"cyberzilla_api_task/model"

	"github.com/labstack/echo/v4"
)

type GetTodoTaskResponse struct {
	Result model.TodoTask `json:"result"`
}

// @Summary Задача
// @Description Возвращает задачу по id
// @Accept  json
// @Produce  json
// @tags TODO
// @in header
// @Param token query string false "Токен списка задач" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Param id path int true "Идентиикатор задачи" default(1)
// @Success 200 {object} GetTodoTaskResponse
// @Failure 400 {object} ErrBadRequest
// @Failure 404 {object} ErrNotFound
// @Failure 500 {object} ErrInternal
// @Router /todos/{id} [get]
func (c *Controller) getTodoTask(ctx echo.Context) error {
	token, err := getAuthorizationToken(ctx)
	if err != nil {
		return c.respondError(ctx, err)
	}
	id, err := getIntParam(ctx, "id")
	if err != nil || id < 1 {
		return c.respondError(ctx, errWrongTodoTaskID)
	}
	res, err := c.app.TodoTask(int64(id), token)
	if err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, res)
}
