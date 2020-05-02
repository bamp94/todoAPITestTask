package controller

import (
	"net/http"

	"cyberzilla_api_task/model"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

// FieldErrors used to returning field validation errors on manual validations with HTTP-code 409 "Conflict"
type FieldErrors map[string]string

// Error mandatory method of error interface
func (f FieldErrors) Error() string {
	var str string
	for key, message := range f {
		str += key + ":" + message + "\n"
	}
	return str
}

// ErrNotFound implements error for HTTP 404
type ErrNotFound struct {
	Msg string
}

// Error implements error interface
func (err ErrNotFound) Error() string {
	if err.Msg == "" {
		return "Запись не найдена"
	}
	return err.Msg
}

func (c *Controller) respondOK(ctx echo.Context, result interface{}) error {
	return ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (c *Controller) respondError(ctx echo.Context, err error) error {
	if fieldErrors, ok := err.(FieldErrors); ok {
		return ctx.JSON(http.StatusBadRequest, gin.H{"fieldErrors": fieldErrors})
	}

	h := gin.H{"error": err.Error()}

	switch {
	case err == model.ErrInternal:
		return ctx.JSON(http.StatusInternalServerError, h)
	case err == model.ErrModelNotFound:
		return ctx.JSON(http.StatusNotFound, h)
	case err.Error() == "multipart: NextPart: EOF":
		return ctx.JSON(http.StatusBadRequest, gin.H{"error": "Тело запроса должно быть в формате multipart/form-data"})
	default:
		return ctx.JSON(http.StatusBadRequest, h)
	}
}
