package controller

import (
	"cyberzilla_api_task/model"

	"github.com/labstack/echo/v4"
)

type TaskDTO struct {
	Task string `json:"task" validate:"required" example:"Do my homework"`
}

// @Summary Создать задачу
// @Description Создает новую задачу
// @Accept  json
// @Produce  json
// @tags Основные
// @in header
// @Param token query string false "Токен списка задач" default(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Param params body TaskDTO true "Тело запроса"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrBadRequest
// @Failure 500 {object} ErrInternal
// @Router /todos [post]
func (c *Controller) createTodoTask(ctx echo.Context) error {
	token, err := getAuthorizationToken(ctx)
	if err != nil {
		return c.respondError(ctx, err)
	}
	var req TaskDTO
	if err := ctx.Bind(&req); err != nil {
		return c.respondError(ctx, err)
	}
	if err := ctx.Validate(&req); err != nil {
		return c.respondError(ctx, err)
	}
	if err := c.app.CreateTodoTask(token, model.TodoTask{
		Task: req.Task,
	}); err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, "ok")
}
