package wgmgmt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type IWgmgmtService interface {
	WgStatus() ServerStatus
	GetPeerInfo() ([]ClientStat, error)
	GetClientConfig(clientID string) (string, string, error)
	GetAllClient() []Client
	GetClient(clientID string) (Client, error)
	AddClient(client Client, autoKeyGen bool) error
	RemoveClient(Client) error
	UpdateClient(Client) error
	SetclientAvailability(clientID string, availability bool) error
	UsedTunelIPS() []string
	GetServer() Server
	AddServer(server Server, autoKeyGen bool) error
	UpdateServer(server Server) error
	ReloadServer() error
	StartServer() error
	StopServer() error
}

type IConfigStorage interface {
	GetServerConfig() Server
	SaveServerConfig(Server) error
	AllClient() []Client
	ClientByID(clientID string) (Client, error)
	ClientTrafficByIP(clientIP string) (SendByte int64, ReceiveByte int64, Error error)
	ClientNumber() (int, int)
	AllTraffic() (SendByte int64, ReceiveByte int64, Error error)
	AddClient(Client) error
	UsedPrivateIPS() []string
	UpdateClient(Client) error
	RemoveClient(Client) error
}

type ILogger interface {
	Infolog(msg string)
	Warnlog(msg string)
	Debuglog(msg string)
	ErrorLog(msg string)
	FatalLog(msg string)
}

type wgService struct {
	configFileNameAndPath string
	interfaceName         string
	configStorage         IConfigStorage
	logger                ILogger
}

//----------------------------------------------------------

func New(cfgFullPath string, interfaceName string, configStorage IConfigStorage, lg ILogger) *wgService {
	return &wgService{
		configFileNameAndPath: cfgFullPath,
		interfaceName:         interfaceName,
		configStorage:         configStorage,
		logger:                lg,
	}
}

//----------------------------------------------------------

func (wg *wgService) WgStatus() ServerStatus {
	isUp := wg.statustWireguard()
	var send, receive int64
	var aErr error
	send, receive, aErr = wg.configStorage.AllTraffic()
	if aErr != nil {
		wg.logger.Warnlog(aErr.Error())
	}
	all, enable := wg.configStorage.ClientNumber()
	status := ServerStatus{
		InterfaceName: wg.interfaceName,
		IsEnable:      isUp,
		AllClient:     all,
		EnableClient:  enable,
		Send:          send,
		Receive:       receive,
	}
	return status
}

// GetPeerInfo :
func (wg *wgService) GetPeerInfo() ([]ClientStat, error) {
	clientsStat := []ClientStat{}
	wgClients := wg.configStorage.AllClient()
	wgInterface, err := wgctrl.New()
	if err != nil {
		return clientsStat, err
	}
	currentdev, errD := wgInterface.Device(wg.interfaceName)
	if errD != nil {
		wg.logger.Warnlog(errD.Error() + "no wireguard interface to read peer data ")
		for _, client := range wgClients {
			var sendByte, receiveByte int64
			var tErr error
			sendByte, receiveByte, tErr = wg.configStorage.ClientTrafficByIP(client.AllocatedIP)
			if tErr != nil {
				wg.logger.Warnlog(tErr.Error())
			}
			clientsStat = append(clientsStat, ClientStat{
				Name:      client.Name,
				PublicKey: client.PublicKey,
				HandShake: HandShake{
					LastHandshakeTime: time.Time{},
					Seen:              false,
				},
				AllocatedIPs: client.AllocatedIP,
				SendByte:     sendByte,
				ReceiveByte:  receiveByte,
			})
		}
		return clientsStat, nil
	}

	for _, peer := range currentdev.Peers {
		handshake := peer.LastHandshakeTime
		seen := handshake.IsZero()
		publicKey := peer.PublicKey.String()
		for _, client := range wgClients {
			if client.PublicKey == publicKey {
				var sendByte, receiveByte int64
				var tErr error
				sendByte, receiveByte, tErr = wg.configStorage.ClientTrafficByIP(client.AllocatedIP)
				if tErr != nil {
					wg.logger.Warnlog(tErr.Error())
				}
				clientsStat = append(clientsStat, ClientStat{
					Name:      client.Name,
					PublicKey: client.PublicKey,
					HandShake: HandShake{
						LastHandshakeTime: handshake,
						Seen:              !seen,
					},
					AllocatedIPs: client.AllocatedIP,
					SendByte:     sendByte,
					ReceiveByte:  receiveByte,
				})
			}
		}

	}
	return clientsStat, nil
}

// GetAllClient :
func (wg *wgService) GetAllClient() []Client {
	return wg.configStorage.AllClient()
}

// GetClientConfig :
func (wg *wgService) GetClientConfig(clientID string) (string, string, error) {
	return wg.createClientConfigByID(clientID)
}

// GetClient :
func (wg *wgService) GetClient(clientID string) (Client, error) {
	return wg.configStorage.ClientByID(clientID)
}

// AddClient :
func (wg *wgService) AddClient(newClient Client, autoKeyGen bool) error {
	newClient.Enabled = true
	srv := wg.configStorage.GetServerConfig()
	_, tunnelNetObj, tErr := net.ParseCIDR(srv.TunnelAddress)
	if tErr != nil {
		return tErr
	}
	clientIpObj := net.ParseIP(newClient.AllocatedIP)
	notInNetwork := tunnelNetObj.Contains(clientIpObj)
	if !notInNetwork {
		return fmt.Errorf("client ip is not in server range.should be in %v", srv.TunnelAddress)
	}
	if autoKeyGen {
		privateKey, priErr := wgtypes.GeneratePrivateKey()
		if priErr != nil {
			return priErr
		}
		publicKey := privateKey.PublicKey().String()
		newClient.PublicKey = publicKey
		newClient.PrivateKey = privateKey.String()
	}
	clientConfigStr, cErr := wg.createClientConfigByValue(newClient)
	if cErr != nil {
		return cErr
	}
	qrdoePNGData, qErr := qrcode.Encode(clientConfigStr, qrcode.Medium, 256)
	if qErr != nil {
		return qErr
	}
	newClient.QRCode = "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(qrdoePNGData))
	err := wg.configStorage.AddClient(newClient)
	if err != nil {
		return err
	}
	sErr := wg.saveConfigToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

// UpdateClient :
func (wg *wgService) UpdateClient(updateClient Client) error {
	uErr := wg.configStorage.UpdateClient(updateClient)
	if uErr != nil {
		return uErr
	}
	clientConfigStr, cErr := wg.createClientConfigByValue(updateClient)
	if cErr != nil {
		return cErr
	}
	qrdoePNGData, qErr := qrcode.Encode(clientConfigStr, qrcode.Medium, 256)
	if qErr != nil {
		return qErr
	}
	updateClient.QRCode = "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(qrdoePNGData))
	sErr := wg.saveConfigToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

// SetclientAvailability :
func (wg *wgService) SetclientAvailability(clientID string, availability bool) error {
	client, err := wg.configStorage.ClientByID(clientID)
	if err != nil {
		return err
	}
	if client.Enabled == availability {
		if client.Enabled {
			return fmt.Errorf("client already enable")
		} else {
			return fmt.Errorf("client already disable")
		}
	}
	client.Enabled = availability
	upErr := wg.configStorage.UpdateClient(client)
	if upErr != nil {
		return upErr
	}
	return nil
}

// RemoveClient :
func (wg *wgService) RemoveClient(rmClient Client) error {
	rErr := wg.configStorage.RemoveClient(rmClient)
	if rErr != nil {
		return rErr
	}
	sErr := wg.saveConfigToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

// GetServer :
func (wg *wgService) GetServer() Server {
	srv := wg.configStorage.GetServerConfig()
	if srv.ListenPort == 0 {
		srv.ListenPort = 51820
	}
	return srv
}

// UpdateServer :
func (wg *wgService) UpdateServer(servercfg Server) error {
	if servercfg.AutoGenerateScript {
		natRule := fmt.Sprintf("iptables -t nat -I POSTROUTING -o %v -j MASQUERADE", servercfg.Interface)
		servercfg.PostUp = []string{natRule}
		natRuleDown := fmt.Sprintf("iptables -t nat -D POSTROUTING -o %v -j MASQUERADE", servercfg.Interface)
		servercfg.PostDown = []string{natRuleDown}
	}
	if servercfg.AutoGenerateKey && servercfg.PrivateKey == "" && servercfg.PublicKey == "" {
		privateKey, priErr := wgtypes.GeneratePrivateKey()
		if priErr != nil {
			return priErr
		}
		publicKey := privateKey.PublicKey().String()
		servercfg.PublicKey = publicKey
		servercfg.PrivateKey = privateKey.String()
	}
	uErr := wg.configStorage.SaveServerConfig(servercfg)
	if uErr != nil {
		return uErr
	}
	sErr := wg.saveConfigToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

// AddServer :
func (wg *wgService) AddServer(newServer Server, autoKeyGen bool) error {
	// PostUp = ufw route allow in on wg0 out on eth0
	// PostUp = iptables -A FORWARD -i wg0 -j ACCEPT
	// PostUp = iptables -t nat -I POSTROUTING -o eth0 -j MASQUERADE
	// PreDown = ufw route delete allow in on wg0 out on eth0
	// PreDown = iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
	if newServer.AutoGenerateScript {
		// allowRule := fmt.Sprintf("iptables -A FORWARD -i %v -j ACCEPT", wg.interfaceName)
		natRule := fmt.Sprintf("iptables -t nat -I POSTROUTING -o %v -j MASQUERADE", newServer.Interface)
		newServer.PostUp = []string{natRule}
		natRuleDown := fmt.Sprintf("iptables -t nat -D POSTROUTING -o %v -j MASQUERADE", newServer.Interface)
		newServer.PostDown = []string{natRuleDown}
	}
	if newServer.AutoGenerateKey {
		privateKey, priErr := wgtypes.GeneratePrivateKey()
		if priErr != nil {
			return priErr
		}
		publicKey := privateKey.PublicKey().String()
		newServer.PublicKey = publicKey
		newServer.PrivateKey = privateKey.String()
	}
	aErr := wg.configStorage.SaveServerConfig(newServer)
	if aErr != nil {
		return nil
	}
	sErr := wg.saveConfigToFile()
	if sErr != nil {
		return sErr
	}
	return nil
}

// UsedTunelIPS :
func (wg *wgService) UsedTunelIPS() []string {
	return wg.configStorage.UsedPrivateIPS()
}

// ReloadServer : reload wireguard by 'wg-quick'
func (wg *wgService) ReloadServer() error {
	wErr := wg.saveConfigToFile()
	if wErr != nil {
		return wErr
	}
	isWgUP := wg.statustWireguard()
	if isWgUP {
		err := wg.stoptWireguard()
		if err != nil {
			return err
		}
	}
	err := wg.startWireguard()
	if err != nil {
		return err
	}
	return nil
}

// StopServer : stop wireguard by 'wg-quick'
func (wg *wgService) StopServer() error {
	isWgUP := wg.statustWireguard()
	if !isWgUP {
		return fmt.Errorf("server is not running")
	}
	err := wg.stoptWireguard()
	if err != nil {
		return err
	}
	return nil

}

// StartServer : start wireguard by 'wg-quick'
func (wg *wgService) StartServer() error {
	isWgUP := wg.statustWireguard()
	if isWgUP {
		return fmt.Errorf("server is running")
	}
	err := wg.startWireguard()
	if err != nil {
		return err
	}
	return nil
}

//-------------------------------------------------
//                Private functions
//-------------------------------------------------

func (wg *wgService) startWireguard() error {
	startCommand := fmt.Sprintf("wg-quick up %v", wg.interfaceName)
	cmd := exec.Command("bash", "-c", startCommand)
	out, err := cmd.CombinedOutput()
	if err != nil && strings.Contains(string(out), "already exists") {
		return fmt.Errorf("alredy running")
	}
	return nil
}

//-------------------------------------------------
func (wg *wgService) stoptWireguard() error {
	stopCommand := fmt.Sprintf("wg-quick down %v", wg.interfaceName)
	cmd := exec.Command("bash", "-c", stopCommand)
	out, err := cmd.CombinedOutput()
	if err != nil && strings.Contains(string(out), "is not a WireGuard interface") {
		return fmt.Errorf("interface not found")
	}
	return nil
}

//-------------------------------------------------
func (wg *wgService) statustWireguard() bool {
	statusCommand := "wg show interfaces"
	cmd := exec.Command("bash", "-c", statusCommand)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmdErr := cmd.Run()
	if cmdErr != nil {
		return false
	}
	if strings.TrimSpace(stdout.String()) == wg.interfaceName {
		return true
	}
	return false

}

//-------------------------------------------------
func (wg *wgService) getTunellIPMask() (int, int) {
	srvConfig := wg.GetServer()
	_, ipNetwork, pErr := net.ParseCIDR(srvConfig.TunnelAddress)
	if pErr != nil {
		return 0, 0
	}
	maskInt, maskbits := ipNetwork.Mask.Size()
	return maskInt, maskbits
}

//-------------------------------------------------
