package store

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"sync"
	"time"
	"wireguard-web-ui/trafficcapture"
	"wireguard-web-ui/usermanagment"
	"wireguard-web-ui/wgmgmt"

	"github.com/google/uuid"
)

type ConfigStore struct {
	dataPath            string
	wgDataFileName      string
	trafficFileName     string
	userAccountFileName string
	wgConfigStore       wgConfig
	logger              ILogger
	internalTableLock   sync.RWMutex
	sourceGroup         map[string]*trafficcapture.SourceGroup
}

type wgConfig struct {
	ServerConfig wgmgmt.Server
	Clients      []wgmgmt.Client
}

type ILogger interface {
	ErrorLog(msg string)
}

//---------------------------------------------
var errClientNotFound = errors.New("client not found")

//---------------------------------------------

// New : create new store for wireguard server and clients config in json format
func New(dataPath string, logger ILogger) (*ConfigStore, error) {
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	wgcfg := wgConfig{
		ServerConfig: wgmgmt.Server{},
		Clients:      []wgmgmt.Client{},
	}
	wgcfg.ServerConfig.PostDown = make([]string, 0)
	wgcfg.ServerConfig.PostUp = make([]string, 0)
	store := ConfigStore{
		dataPath:            dataPath,
		wgDataFileName:      path.Join(dataPath, "wgData.json"),
		trafficFileName:     path.Join(dataPath, "traffic.json"),
		userAccountFileName: path.Join(dataPath, "userAccount.json"),
		wgConfigStore:       wgcfg,
		sourceGroup:         make(map[string]*trafficcapture.SourceGroup),
		logger:              logger,
	}
	//---------------
	_, stErr := os.Stat(store.wgDataFileName)
	if errors.Is(stErr, os.ErrNotExist) {
	} else {
		wgFileData, err := ioutil.ReadFile(store.wgDataFileName)
		if err != nil {
			return nil, err
		}
		if len(wgFileData) != 0 {
			err = json.Unmarshal(wgFileData, &store.wgConfigStore)
			if err != nil {
				return nil, err
			}
		}
	}

	//------------------
	_, stErrt := os.Stat(store.trafficFileName)
	if errors.Is(stErrt, os.ErrNotExist) {
	} else {
		fileData, err := ioutil.ReadFile(store.trafficFileName)
		if err != nil {
			return nil, err
		}
		if len(fileData) != 0 {
			err = json.Unmarshal(fileData, &store.sourceGroup)
			if err != nil {
				return nil, err
			}
		}
	}

	//------------------
	_, useFStatErr := os.Stat(store.userAccountFileName)
	if errors.Is(useFStatErr, os.ErrNotExist) {
		user := "wgadmin"
		clearPassword := "wireguardadmin"
		md5Hasher := md5.New()
		md5Hasher.Write([]byte(clearPassword))
		md5Password := hex.EncodeToString(md5Hasher.Sum(nil))

		admin := usermanagment.AdminUser{
			ID:       uuid.New().String(),
			UserName: user,
			Password: md5Password,
		}
		userList := []usermanagment.AdminUser{admin}

		adminListByte, mErr := json.Marshal(userList)
		if mErr != nil {
			return nil, mErr
		}
		wErr := ioutil.WriteFile(store.userAccountFileName, adminListByte, 0666)
		if wErr != nil {
			return nil, wErr
		}
	}
	//------------------
	return &store, nil
}

//---------------------------------------------

func (cs *ConfigStore) GetServerConfig() wgmgmt.Server {
	return cs.wgConfigStore.ServerConfig
}

//---------------------------------------------
func (cs *ConfigStore) SaveServerConfig(serverConf wgmgmt.Server) error {
	cs.wgConfigStore.ServerConfig = serverConf
	sErr := cs.saveWgDataToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

//---------------------------------------------
func (cs *ConfigStore) AllClient() []wgmgmt.Client {
	cs.internalTableLock.Lock()
	defer cs.internalTableLock.Unlock()
	clients := cs.wgConfigStore.Clients
	return clients
}

//---------------------------------------------
func (cs *ConfigStore) ClientByID(clientID int64) (wgmgmt.Client, error) {
	sClient := wgmgmt.Client{}
	for _, client := range cs.wgConfigStore.Clients {
		if client.ID == clientID {
			sClient := wgmgmt.Client{}
			sClient = client
			return sClient, nil
		}
	}
	return sClient, errClientNotFound
}

//---------------------------------------------
func (cs *ConfigStore) AddClient(newClient wgmgmt.Client) error {
	clientID := time.Now().UnixNano()
	newClient.ID = clientID

	cs.wgConfigStore.Clients = append(cs.wgConfigStore.Clients, newClient)

	sErr := cs.saveWgDataToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

//---------------------------------------------
func (cs *ConfigStore) UpdateClient(updateClient wgmgmt.Client) error {
	found := false

	for index, client := range cs.wgConfigStore.Clients {
		if client.ID == updateClient.ID {
			cs.wgConfigStore.Clients[index] = updateClient
			found = true
		}
	}
	if !found {
		return errClientNotFound
	}
	sErr := cs.saveWgDataToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

//---------------------------------------------
func (cs *ConfigStore) RemoveClient(rmClient wgmgmt.Client) error {
	found := false
	for index, client := range cs.wgConfigStore.Clients {
		if client.ID == rmClient.ID {
			cs.wgConfigStore.Clients = append(cs.wgConfigStore.Clients[:index], cs.wgConfigStore.Clients[index+1:]...)
			found = true
		}
	}
	if !found {
		return errClientNotFound
	}
	sErr := cs.saveWgDataToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

//---------------------------------------------
func (cs *ConfigStore) AvailablePrivateIPS() []string {
	srvConfig := cs.wgConfigStore.ServerConfig
	serverIPRange := srvConfig.TunnelAddress
	iplist, _, err := cidrToIpList(serverIPRange)
	if err != nil {
		cs.logger.ErrorLog(err.Error())
		return []string{}
	}
	ipSrv, _, pErr := net.ParseCIDR(serverIPRange)
	if pErr != nil {
		cs.logger.ErrorLog(pErr.Error())
		return []string{}
	}
	iplist = removeFromStringArr(iplist, ipSrv.String())
	allClient := cs.wgConfigStore.Clients
	for _, client := range allClient {
		iplist = removeFromStringArr(iplist, client.AllocatedIP)
	}
	return iplist
}

//---------------------------------------------
func (cs *ConfigStore) UsedPrivateIPS() []string {
	usedIp := []string{}
	srvConfig := cs.wgConfigStore.ServerConfig
	serverIPRange := srvConfig.TunnelAddress
	ipSrv, ipnet, pErr := net.ParseCIDR(serverIPRange)
	if pErr != nil {
		cs.logger.ErrorLog(pErr.Error())
		return usedIp
	}
	netID := ipnet.IP.String()
	usedIp = append(usedIp, ipSrv.String())
	usedIp = append(usedIp, netID)
	broadcast := net.IP(net.ParseIP("0.0.0.0").To4())
	for i := 0; i < len(ipnet.IP); i++ {
		broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
	}
	usedIp = append(usedIp, broadcast.String())
	allClient := cs.wgConfigStore.Clients
	for _, client := range allClient {
		usedIp = append(usedIp, client.AllocatedIP)
	}
	return usedIp
}

//---------------------------------------------
func (cs *ConfigStore) WriteSourceGroup(sourceIps map[string]*trafficcapture.SourceGroup) {
	for ip := range sourceIps {
		_, found := cs.sourceGroup[ip]
		if found {
			cs.sourceGroup[ip].PacketCount += sourceIps[ip].PacketCount
			cs.sourceGroup[ip].RecieveByte += sourceIps[ip].RecieveByte
			cs.sourceGroup[ip].SendByte += sourceIps[ip].SendByte
		} else {
			cs.sourceGroup[ip] = sourceIps[ip]
		}
	}
	err := cs.saveTrafficToFile()
	if err != nil {
		log.Println(err)
	}
}

//---------------------------------------------
func (cs *ConfigStore) ClientTrafficByIP(clientIP string) (int64, int64, error) {
	for ip := range cs.sourceGroup {
		if ip == clientIP {
			return cs.sourceGroup[ip].SendByte, cs.sourceGroup[ip].RecieveByte, nil
		}
	}
	return 0, 0, nil
}

//---------------------------------------------
func (cs *ConfigStore) AllTraffic() (int64, int64, error) {
	var send, receive int64
	for _, info := range cs.sourceGroup {
		send += info.SendByte
		receive += info.RecieveByte
	}
	return send, receive, nil
}

//---------------------------------------------
func (cs *ConfigStore) ClientNumber() (int, int) {
	cs.internalTableLock.Lock()
	defer cs.internalTableLock.Unlock()
	clients := cs.wgConfigStore.Clients
	var allClient, enableClient int
	for _, client := range clients {
		allClient++
		if client.Enabled {
			enableClient++
		}
	}
	return allClient, enableClient
}
func (cs *ConfigStore) CheckUserAndPassword(adminUser usermanagment.AdminUser) bool {
	userList := []usermanagment.AdminUser{}
	wgFileData, err := ioutil.ReadFile(cs.userAccountFileName)
	if err != nil {
		cs.logger.ErrorLog(err.Error())
		// fmt.Println(err)
		return false
	}
	err = json.Unmarshal(wgFileData, &userList)
	if err != nil {
		cs.logger.ErrorLog(err.Error())
		// fmt.Println(err)
		return false
	}
	userIsValid := false
	clearPassword := adminUser.Password
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(clearPassword))
	md5Password := hex.EncodeToString(md5Hasher.Sum(nil))
	for _, user := range userList {
		if user.UserName == adminUser.UserName && user.Password == md5Password {
			userIsValid = true
		}
	}
	return userIsValid
}

//***********************************************//
func (cs *ConfigStore) saveWgDataToFile() error {
	configByte, mErr := json.Marshal(cs.wgConfigStore)
	if mErr != nil {
		return mErr
	}
	wErr := ioutil.WriteFile(cs.wgDataFileName, configByte, 0666)
	if wErr != nil {
		return wErr
	}
	return nil
}

func (cs *ConfigStore) saveTrafficToFile() error {
	configByte, mErr := json.Marshal(cs.sourceGroup)
	if mErr != nil {
		return mErr
	}
	wErr := ioutil.WriteFile(cs.trafficFileName, configByte, 0666)
	if wErr != nil {
		return wErr
	}
	return nil
}
