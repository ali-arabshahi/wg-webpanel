package http

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func (httpsrv httpServer) jsonResponse(ctx echo.Context, code int, data interface{}, err error) error {
	msg := ""
	if err != nil {
		ip := ""
		ippart := strings.Split(strings.TrimSpace(ctx.Request().RemoteAddr), ":")
		if len(ippart) > 1 {
			ip = ippart[0]
		}
		httpsrv.routerloger.HttpErrorLog(ctx.Request().RequestURI, ctx.Request().Method, ip, code, err.Error())
		msg = err.Error()
	} else {
		msg = "ok"
	}
	response := struct {
		Code    int
		Data    interface{}
		Message string
	}{
		code,
		data,
		msg,
	}

	return ctx.JSONPretty(code, response, "  ")
}
