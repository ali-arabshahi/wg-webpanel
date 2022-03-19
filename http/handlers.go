package http

import (
	"fmt"
	"net/http"
	"time"
	"wireguard-web-ui/system"
	"wireguard-web-ui/usermanagment"
	"wireguard-web-ui/wgmgmt"

	"github.com/labstack/echo/v4"
)

//---------------------------------------------------------------------------

func (httpsrv httpServer) login(userMgmg usermanagment.IUserManagment) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		userCred := usermanagment.AdminUser{}
		err := ctx.Bind(&userCred)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", err)
		}
		if userCred.UserName == "" || userCred.Password == "" {
			return httpsrv.jsonResponse(ctx, http.StatusNotAcceptable, "user or pass is missing", fmt.Errorf("user or pass is missing"))
		}
		cErr := userMgmg.CheckUserInfo(userCred)
		if cErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusUnauthorized, nil, cErr)
		}
		token, tErr := userMgmg.GenerateToken(userCred.UserName)
		if tErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusUnauthorized, nil, tErr)
		}
		res := struct {
			User  string `json:"user"`
			Token string `json:"token"`
		}{
			userCred.UserName,
			token,
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, res, nil)
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) getServerConfig(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		srvConfig := wg.GetServer()
		return httpsrv.jsonResponse(ctx, http.StatusOK, srvConfig, nil)
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) saveServerConfig(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		srv := wgmgmt.Server{}
		err := ctx.Bind(&srv)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", err)
		}
		vErr := modelValidator.Struct(srv)
		if vErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", vErr)
		}
		sErr := wg.UpdateServer(srv)
		if sErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", sErr)
		}
		time.Sleep(1 * time.Second)
		return ctx.JSONPretty(200, "ok", "  ")
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) interfaces(systemService system.Isystem) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		nics, err := systemService.NetworkInterfaces()
		if err != nil {
			httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", err)
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, nics, nil)
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) getAllClients(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		clients := wg.GetAllClient()
		return httpsrv.jsonResponse(ctx, http.StatusOK, clients, nil)

	}
}

//---------------------------------------------------------------------------

// addClient : add client
func (httpsrv httpServer) addClient(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		client := wgmgmt.Client{}
		err := ctx.Bind(&client)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", err)
		}
		vErr := modelValidator.Struct(client)
		if vErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", vErr)
		}
		sErr := wg.AddClient(client, true)
		if sErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", sErr)
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, nil, nil)

	}
}

//---------------------------------------------------------------------------

// updateClient : add client
func (httpsrv httpServer) updateClient(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		client := wgmgmt.Client{}
		err := ctx.Bind(&client)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", err)
		}
		sErr := wg.UpdateClient(client)
		if sErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", sErr)
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, nil, nil)

	}
}

//---------------------------------------------------------------------------
func (httpsrv httpServer) setclientAvailability(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		client := wgmgmt.Client{}
		err := ctx.Bind(&client)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", err)
		}
		sErr := wg.SetclientAvailability(client.ID, client.Enabled)
		if sErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", sErr)
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, nil, nil)

	}
}

//---------------------------------------------------------------------------
// deleteClient : add client
func (httpsrv httpServer) deleteClient(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		client := wgmgmt.Client{}
		err := ctx.Bind(&client)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", err)
		}
		sErr := wg.RemoveClient(client)
		if sErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", sErr)
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, nil, nil)

	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) getUsedTunelIps(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		usedIPS := wg.UsedTunelIPS()
		return httpsrv.jsonResponse(ctx, http.StatusOK, usedIPS, nil)
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) getClientConfig(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		clientID := ctx.Param("id")
		if clientID == "" {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", fmt.Errorf("need parameter"))
		}
		config, fname, err := wg.GetClientConfig(clientID)
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", err)
		}

		data := []byte(config)
		ctx.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", "attachment", "client-config.txt"))
		ctx.Response().Header().Set("X-FileName", fname+"-cfg.txt")
		return ctx.Blob(http.StatusOK, "text/plain", data)
		// return ctx.String(200, config)
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) reload(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		err := wg.ReloadServer()
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", err)
		}
		return ctx.JSONPretty(200, "reload successful", "  ")
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) startStop(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		operation := ctx.Param("operation")
		if operation == "" {
			return httpsrv.jsonResponse(ctx, http.StatusBadRequest, "", fmt.Errorf("need parameter"))
		}
		var actionErr error
		if operation == "start" {
			actionErr = wg.StartServer()
		} else if operation == "stop" {
			actionErr = wg.StopServer()
		} else {
			actionErr = fmt.Errorf("unknown action")
			return httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", actionErr)
		}
		if actionErr != nil {
			return httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", actionErr)
		}
		return ctx.JSONPretty(200, "ok", "  ")
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) getClientStat(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		stat, err := wg.GetPeerInfo()
		if err != nil {
			return httpsrv.jsonResponse(ctx, http.StatusInternalServerError, "", err)
		}
		return httpsrv.jsonResponse(ctx, http.StatusOK, stat, nil)
	}
}

//---------------------------------------------------------------------------

func (httpsrv httpServer) serverStatus(wg wgmgmt.IWgmgmtService) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		status := wg.WgStatus()
		return httpsrv.jsonResponse(ctx, http.StatusOK, status, nil)

	}
}
