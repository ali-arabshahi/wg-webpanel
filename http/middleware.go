package http

import (
	"fmt"
	"net/http"
	"strings"
	"wireguard-web-ui/usermanagment"

	"github.com/labstack/echo/v4"
)

// CheckToenvalidation : check token validity
func (httpsrv httpServer) checkToenvalidation(userMgmg usermanagment.IUserManagment) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token := ctx.Request().Header.Get("token")
			if token == "" {
				return httpsrv.jsonResponse(ctx, http.StatusUnauthorized, "token not found", fmt.Errorf("token not found"))

			}
			_, validErr := userMgmg.ValidateToken(token)
			if validErr != nil {
				return httpsrv.jsonResponse(ctx, http.StatusUnauthorized, "token is not valid", fmt.Errorf("token is not valid"))
			}
			return next(ctx)
		}
	}
}

// accessLogMidelware : access log midelware
func (httpsrv httpServer) accessLogMidelware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ip := ""
		ippart := strings.Split(strings.TrimSpace(ctx.Request().RemoteAddr), ":")
		if len(ippart) > 1 {
			ip = ippart[0]
		}
		httpsrv.routerloger.AccessLog(ctx.Request().RequestURI, ctx.Request().Method, ip)
		return next(ctx)
	}
}

