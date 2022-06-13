package main

import (
	"flag"
	"log"
	"math/rand"
	"strings"
	"time"
	"wireguard-web-ui/http"
	"wireguard-web-ui/serverconfig"
	"wireguard-web-ui/servicelog"
	"wireguard-web-ui/store"
	"wireguard-web-ui/system"
	"wireguard-web-ui/trafficaggrigate"
	"wireguard-web-ui/trafficcapture"
	"wireguard-web-ui/usermanagment"
	"wireguard-web-ui/wgmgmt"
)

func main() {
	//------------------------------
	configPath := flag.String("config", "./config.json", "server config file path")
	flag.Parse()
	//------------------------------
	config, cErr := serverconfig.LoadConfig(*configPath)
	if cErr != nil {
		log.Fatalln(cErr)
	}
	//-----------------------------
	loggerFileConfig := servicelog.LogFileConfig{
		FileAddress: config.LogFileAddress,
		MaxSize:     20,
		MaxBackups:  2,
		MaxAge:      30,
	}
	logger := servicelog.New("debug", &loggerFileConfig, true)
	//-----------------------------
	store, stErr := store.NewStore(store.SqlliteType, config.DataDirectory, logger)
	if stErr != nil {
		log.Fatalln(stErr)
	}
	//-----------------------------
	parts := strings.Split(config.WireguardConfigPath, "/")
	fName := parts[len(parts)-1]
	interfaceName := strings.Split(fName, ".")[0]
	//-----------------------------
	systemService := system.New()
	//-----------------------------
	wgService := wgmgmt.New(config.WireguardConfigPath, interfaceName, store, logger)
	//-----------------------------
	analyzer, err := trafficcapture.New(interfaceName, true, 65000, logger)
	if err != nil {
		log.Fatalln(err)
	}
	//-----------------------------
	sourceAggrigator, errG := trafficaggrigate.New(trafficaggrigate.SourceGroup, 5, store)
	if errG != nil {
		log.Println(errG)
	}
	//-----------------------------
	analyzer.StartCapturing()
	analyzer.AddProcessor("source-group", sourceAggrigator)
	//-----------------------------
	tokenSecretKey := randomString(12)
	userMgmtService := usermanagment.New(tokenSecretKey, store)
	//-----------------------------
	httpServerConfig := http.RouterConfig{
		Port:      config.ServerPort,
		StaticDir: config.StaticDir,
		Cert:      config.CertAddress,
		Certkey:   config.CertKeyAddress,
	}
	http.RunRouter(httpServerConfig, config.EnableHTTPS, wgService, systemService, userMgmtService, logger)
}

//************************************************************************
//************************************************************************

func randomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
