package controller

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type HealthCheck struct {
	Result struct {
		DB string `json:"DB,omitempty" example:"ok"`
	} `json:"result"`
}

// @Summary Проверка сервера
// @Description Запускает процесс проверки всех зависимостей, необходимых для корректной работы сервера
// @Accept  json
// @Produce  json
// @tags Служебные
// @Success 200 {object} HealthCheck
// @Router /healthcheck [get]
func (c *Controller) healthcheck(ctx echo.Context) error {
	healthcheck := make(map[string]string)
	mx := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := c.app.PingDatabase()

		mx.Lock()
		if err != nil {
			healthcheck["DB"] = err.Error()
		} else {
			healthcheck["DB"] = "ok"
		}
		mx.Unlock()
		wg.Done()
	}()

	wg.Wait()
	return ctx.JSON(http.StatusOK, echo.Map{"result": healthcheck})
}
