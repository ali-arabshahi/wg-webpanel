package wgmgmt

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"text/template"
)

const wgConfigTemplate = `[Interface]
Address = {{.Server.TunnelAddress}}
{{- range .Server.PostUp}}
PostUp = {{ . }}
{{- end}}
{{- range .Server.PostDown}}
PreDown = {{ . }}
{{- end}}
ListenPort = {{.Server.ListenPort}}
PrivateKey = {{.Server.PrivateKey}}

# - - - - - Clients Section - - - - -
{{- range .Clients}}{{if eq .Enabled true}}
# ID:    {{ .ID }}
# Name:  {{ .Name }}
[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = {{ .AllocatedIpCIDR}}
{{end}}
{{end}}
`

// AllowedIPs = {{ JoinStr .AllowedIPs "," -}}
const clientConfigTemplate = `[Interface]
PrivateKey = {{.ClientPriKey}}
Address = {{.ClientAddress}}
DNS = {{.ClientDNS}}

[Peer]
PublicKey = {{.ServerPublicKey}}
Endpoint = {{.Endpoint}}
AllowedIPs = {{ JoinStr .ClientAllowIPS "," -}}
`

// {{ JoinStr .PostScript "," }}

func (wg *wgService) saveConfigToFile() error {
	// create text template
	configTemplate := template.New("wg-config")
	configTemplate = configTemplate.Funcs(template.FuncMap{"JoinStr": strings.Join})
	configTemplate, err := configTemplate.Parse(wgConfigTemplate)
	if err != nil {
		return err
	}
	_, bits := wg.getTunellIPMask()
	serverConfig := wg.GetServer()
	clientsConfig := wg.GetAllClient()
	for id, client := range clientsConfig {
		if client.Enabled {
			clientsConfig[id].AllocatedIpCIDR = fmt.Sprintf("%s/%v", clientsConfig[id].AllocatedIP, bits)
		}
	}
	srvConfig := struct {
		Server  Server
		Clients []Client
	}{
		serverConfig,
		clientsConfig,
	}
	//--------------------------------------------
	// create or truncate config file to wtite data
	configFile, err := os.Create(wg.configFileNameAndPath)
	if err != nil {
		return err
	}
	defer configFile.Close()
	//--------------------------------------------
	// insert data in template and write to file
	exErr := configTemplate.Execute(configFile, &srvConfig)
	if exErr != nil {
		return exErr
	}
	return nil
}

func (wg *wgService) createClientConfigByID(clientID string) (string, string, error) {
	clientTemplate := template.New("client-config")
	clientTemplate = clientTemplate.Funcs(template.FuncMap{"JoinStr": strings.Join})
	clientTemplate, err := clientTemplate.Parse(clientConfigTemplate)
	if err != nil {
		return "", "", err
	}
	//------ server config
	srvConfig := wg.GetServer()
	_, _, pErr := net.ParseCIDR(srvConfig.TunnelAddress)
	if pErr != nil {
		return "", "", pErr
	}
	mask, _ := wg.getTunellIPMask()
	maskStr := fmt.Sprint(mask)
	srvEndPoint := fmt.Sprintf("%v:%v", srvConfig.Address, srvConfig.ListenPort)
	//------ client config
	client, cErr := wg.configStorage.ClientByID(clientID)
	if cErr != nil {
		return "", "", cErr
	}
	clientIP := fmt.Sprintf("%v/%v", client.AllocatedIP, maskStr)
	clientConfig := struct {
		ClientAddress   string
		ClientPriKey    string
		ClientDNS       string
		ClientAllowIPS  []string
		ServerPublicKey string
		Endpoint        string
	}{
		clientIP,
		client.PrivateKey,
		client.DNSAddress,
		client.AllowedIPs,
		srvConfig.PublicKey,
		srvEndPoint,
	}
	//--------------------------------------------
	// write template with data injected
	var clientConfBuf bytes.Buffer
	exErr := clientTemplate.Execute(&clientConfBuf, &clientConfig)
	if exErr != nil {
		return "", "", exErr
	}
	return clientConfBuf.String(), client.Name, nil
}

func (wg *wgService) createClientConfigByValue(client Client) (string, error) {
	clientTemplate := template.New("client-config")
	clientTemplate = clientTemplate.Funcs(template.FuncMap{"JoinStr": strings.Join})
	clientTemplate, err := clientTemplate.Parse(clientConfigTemplate)
	if err != nil {
		return "", err
	}
	//------ server config
	srvConfig := wg.GetServer()
	_, ipNetwork, pErr := net.ParseCIDR(srvConfig.TunnelAddress)
	if pErr != nil {
		return "", pErr
	}
	maskInt, _ := ipNetwork.Mask.Size()
	maskStr := fmt.Sprint(maskInt)
	srvEndPoint := fmt.Sprintf("%v:%v", srvConfig.Address, srvConfig.ListenPort)
	//------ client config
	clientIP := fmt.Sprintf("%v/%v", client.AllocatedIP, maskStr)
	clientConfig := struct {
		ClientAddress   string
		ClientPriKey    string
		ClientDNS       string
		ClientAllowIPS  []string
		ServerPublicKey string
		Endpoint        string
	}{
		clientIP,
		client.PrivateKey,
		client.DNSAddress,
		client.AllowedIPs,
		srvConfig.PublicKey,
		srvEndPoint,
	}
	//--------------------------------------------
	// write template with data injected
	var clientConfBuf bytes.Buffer
	exErr := clientTemplate.Execute(&clientConfBuf, &clientConfig)
	if exErr != nil {
		return "", exErr
	}
	return clientConfBuf.String(), nil
}
