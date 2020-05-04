package controller

import (
	"github.com/labstack/echo/v4"
)

type ProxyDTO struct {
	ProxyServers []string `json:"proxyServers" validate:"required" example:"192.168.1.1:80,172.24.130.50:256,5.196.246.54:8080"`
}

// @Summary Проверить прокси
// @Description Проверяет список прокси серверов и возвращает их статус
// @Accept  json
// @Produce  json
// @tags Proxy
// @in header
// @Param params body ProxyDTO true "Тело запроса"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrBadRequest
// @Failure 500 {object} ErrInternal
// @Router /check [post]
func (c *Controller) checkProxyServers(ctx echo.Context) error {
	var req ProxyDTO
	if err := ctx.Bind(&req); err != nil {
		return c.respondError(ctx, err)
	}
	if err := ctx.Validate(&req); err != nil {
		return c.respondError(ctx, err)
	}
	return c.respondOK(ctx, c.app.ProxyServersStatuses(req.ProxyServers))
}
