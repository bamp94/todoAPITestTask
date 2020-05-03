package controller

import (
	"github.com/labstack/echo/v4"
)

// @Summary Удалить задачу
// @Description Удаляет задачу по id
// @Accept  json
// @Produce  json
// @tags Основные
// @in header
// @Param token query string false "Токен списка задач" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Param id path int true "Идентиикатор задачи" default(1)
// @Success 200 {object} GetTodoTaskResponse
// @Failure 400 {object} ErrBadRequest
// @Failure 404 {object} ErrNotFound
// @Failure 500 {object} ErrInternal
// @Router /todos/{id} [delete]
func (c *Controller) deleteTodoTask(ctx echo.Context) error {
	token, err := getAuthorizationToken(ctx)
	if err != nil {
		return c.respondError(ctx, err)
	}
	id, err := getIntParam(ctx, "id")
	if err != nil || id < 1 {
		return c.respondError(ctx, errWrongTodoTaskID)
	}
	if err := c.app.DeleteTodoTask(int64(id), token); err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, "ok")
}
