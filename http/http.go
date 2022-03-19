package http

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"wireguard-web-ui/system"
	"wireguard-web-ui/usermanagment"
	"wireguard-web-ui/wgmgmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterConfig struct {
	Port      string
	StaticDir string
	Cert      string
	Certkey   string
}

type ILogger interface {
	AccessLog(Url string, Method string, RemoteAddr string)
	HttpErrorLog(Url string, Method string, RemoteAddr string, HttpCode int, Details string)
}

type httpServer struct {
	router       *echo.Echo
	serverConfig RouterConfig
	routerloger  ILogger
}

// RunRouter : register route and run http server
func RunRouter(serverCfg RouterConfig, enableSSL bool, wgService wgmgmt.IWgmgmtService, systemService system.Isystem, userMgmtService usermanagment.IUserManagment, logger ILogger) error {
	httpSrv := httpServer{
		router:       echo.New(),
		serverConfig: serverCfg,
		routerloger:  logger,
	}
	newValidator()
	httpSrv.router = echo.New()
	// cors options -------------------------------------------
	httpSrv.router.Use(httpSrv.accessLogMidelware)
	httpSrv.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST,GET,PUT,DELETE,PATCH"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "token"},
		ExposeHeaders:    []string{"X-FileName", "Content-Disposition"},
		AllowCredentials: true,
	}))
	// set static config -------------------------------------------
	httpSrv.router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   httpSrv.serverConfig.StaticDir,
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))
	// rest API endpoints ------------------------------------------
	httpSrv.router.POST("/login", httpSrv.login(userMgmtService))
	private := httpSrv.router.Group("/private", httpSrv.checkToenvalidation(userMgmtService))
	private.GET("/server", httpSrv.getServerConfig(wgService))
	private.POST("/server", httpSrv.saveServerConfig(wgService))
	private.POST("/server/reload", httpSrv.reload(wgService))
	private.POST("/server/:operation", httpSrv.startStop(wgService))
	private.GET("/server/status", httpSrv.serverStatus(wgService))
	//--------------------------------------------------------
	private.GET("/interfaces", httpSrv.interfaces(systemService))
	//--------------------------------------------------------
	private.GET("/client", httpSrv.getAllClients(wgService))
	private.POST("/client", httpSrv.addClient(wgService))
	private.PUT("/client", httpSrv.updateClient(wgService))
	private.DELETE("/client", httpSrv.deleteClient(wgService))
	private.GET("/client/stat", httpSrv.getClientStat(wgService))
	private.GET("/client/config/:id", httpSrv.getClientConfig(wgService))
	private.GET("/client/usedIP", httpSrv.getUsedTunelIps(wgService))
	private.PUT("/client/availability", httpSrv.setclientAvailability(wgService))
	//http server config -------------------------------------
	server := &http.Server{
		Addr:         "0.0.0.0:" + httpSrv.serverConfig.Port,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	if enableSSL {
		cer, err := tls.LoadX509KeyPair(httpSrv.serverConfig.Cert, httpSrv.serverConfig.Certkey)
		if err != nil {
			log.Println(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}
		server.TLSConfig = config
	}
	httpSrv.router.Logger.Fatal(httpSrv.router.StartServer(server))
	//--------------------------------------------------------
	return nil
}
