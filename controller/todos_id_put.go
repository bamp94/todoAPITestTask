package controller

import (
	"cyberzilla_api_task/model"

	"github.com/labstack/echo/v4"
)

// @Summary Обновить задачу
// @Description Обновляет данные задачи по id
// @Accept  json
// @Produce  json
// @tags TODO
// @in header
// @Param token query string false "Токен списка задач" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Param id path int true "Идентиикатор задачи" default(1)
// @Param params body TaskDTO true "Тело запроса"
// @Success 200 {object} GetTodoTaskResponse
// @Failure 400 {object} ErrBadRequest
// @Failure 404 {object} ErrNotFound
// @Failure 500 {object} ErrInternal
// @Router /todos/{id} [put]
func (c *Controller) updateTodoTask(ctx echo.Context) error {
	token, err := getAuthorizationToken(ctx)
	if err != nil {
		return c.respondError(ctx, err)
	}
	id, err := getIntParam(ctx, "id")
	if err != nil || id < 1 {
		return c.respondError(ctx, errWrongTodoTaskID)
	}
	var req TaskDTO
	if err := ctx.Bind(&req); err != nil {
		return c.respondError(ctx, err)
	}
	if err := ctx.Validate(&req); err != nil {
		return c.respondError(ctx, err)
	}
	task := model.TodoTask{
		ID:   int64(id),
		Task: req.Task,
	}
	if err := c.app.UpdateTodoTask(token, task); err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, "ok")
}
