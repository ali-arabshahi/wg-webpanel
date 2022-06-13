package store

import (
	"fmt"
	"path"
	"reflect"
	"testing"
	"time"
	"wireguard-web-ui/wgmgmt"
)

type mockLogger struct{}

func (ml mockLogger) ErrorLog(msg string) { fmt.Println(msg) }

func TestNew(t *testing.T) {
	tmp := t.TempDir()
	tmp2 := t.TempDir()
	fAdress := path.Join(tmp)
	fAdress2 := path.Join(tmp2)
	mockLogger := mockLogger{}
	store, err := newFileStore(fAdress, mockLogger)
	if err != nil {
		t.Errorf("New() = %v", err)
		return
	}
	store2, err := newFileStore(fAdress2, mockLogger)
	if err != nil {
		t.Errorf("New() = %v", err)
		return
	}
	type args struct {
		fileAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    *fileStore
		wantErr bool
	}{
		{"file not exist",
			args{
				fileAddr: fAdress,
			},
			store,
			false,
		},
		{"file exist",
			args{
				fileAddr: fAdress2,
			},
			store2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFileStore(tt.args.fileAddr, mockLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigStore_GetServerConfig(t *testing.T) {

	tests := []struct {
		name string
		cs   *fileStore
		want wgmgmt.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.GetServerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigStore.GetServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigStore_SaveServerConfig(t *testing.T) {
	mockLogger := mockLogger{}
	tmp := t.TempDir()
	fAdress := path.Join(tmp)
	store, err := newFileStore(fAdress, mockLogger)
	if err != nil {
		t.Errorf("New() = %v", err)
		return
	}
	type args struct {
		serverConf wgmgmt.Server
	}
	sampleServer := wgmgmt.Server{
		Address:    "192.168.10.1",
		ListenPort: 514,
		PostUp:     []string{"post script-1", "post script-2"},
		PostDown:   []string{"pre script"},
		PrivateKey: "pri key",
		PublicKey:  "pub key",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	tests := []struct {
		name    string
		cs      *fileStore
		args    args
		wantErr bool
	}{
		{
			"success",
			store,
			args{
				serverConf: sampleServer,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.SaveServerConfig(tt.args.serverConf); (err != nil) != tt.wantErr {
				t.Errorf("ConfigStore.SaveServerConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigStore_AllClient(t *testing.T) {
	tests := []struct {
		name string
		cs   *fileStore
		want []wgmgmt.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.AllClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigStore.AllClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigStore_ClientByID(t *testing.T) {
	tmp := t.TempDir()
	mockLogger := mockLogger{}
	fAdress := path.Join(tmp)
	store, err := newFileStore(fAdress, mockLogger)
	if err != nil {
		t.Errorf("New() = %v", err)
		return
	}
	sampleClient := wgmgmt.Client{
		PrivateKey:  "pri key",
		PublicKey:   "pub key",
		Name:        "ali",
		AllocatedIP: "192.168.10.2/32",
		AllowedIPs:  []string{"0.0.0.0/0"},
		QRCode:      "qcode",
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	inErr := store.AddClient(sampleClient)
	if inErr != nil {
		t.Errorf("AddClient() = %v", inErr)
		return
	}
	clients := store.AllClient()
	getclient := clients[0]
	type args struct {
		clientID int64
	}
	tests := []struct {
		name    string
		cs      *fileStore
		args    args
		want    wgmgmt.Client
		wantErr bool
	}{
		{
			"found",
			store,
			args{
				clientID: getclient.ID,
			},
			getclient,
			false,
		},
		{
			"not found",
			store,
			args{
				clientID: 11111,
			},
			wgmgmt.Client{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cs.ClientByID(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigStore.ClientByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigStore.ClientByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigStore_AddClient(t *testing.T) {
	tmp := t.TempDir()
	fAdress := path.Join(tmp)
	mockLogger := mockLogger{}
	store, err := newFileStore(fAdress, mockLogger)
	if err != nil {
		t.Errorf("New() = %v", err)
		return
	}
	type args struct {
		newClient wgmgmt.Client
	}
	cl1ID := time.Now().UnixNano()
	sampleClient := wgmgmt.Client{
		ID:          cl1ID,
		PrivateKey:  "pri key",
		PublicKey:   "pub key",
		Name:        "ali",
		AllocatedIP: "192.168.10.2/32",
		AllowedIPs:  []string{"0.0.0.0/0"},
		QRCode:      "qcode",
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tests := []struct {
		name    string
		cs      *fileStore
		args    args
		wantErr bool
	}{
		{
			"success",
			store,
			args{
				newClient: sampleClient,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.AddClient(tt.args.newClient); (err != nil) != tt.wantErr {
				t.Errorf("ConfigStore.AddClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigStore_UpdateClient(t *testing.T) {
	type args struct {
		updateClient wgmgmt.Client
	}
	tests := []struct {
		name    string
		cs      *fileStore
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.UpdateClient(tt.args.updateClient); (err != nil) != tt.wantErr {
				t.Errorf("ConfigStore.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigStore_RemoveClient(t *testing.T) {
	type args struct {
		rmClient wgmgmt.Client
	}
	tests := []struct {
		name    string
		cs      *fileStore
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.RemoveClient(tt.args.rmClient); (err != nil) != tt.wantErr {
				t.Errorf("ConfigStore.RemoveClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigStore_saveToFile(t *testing.T) {
	mockLogger := mockLogger{}
	tmp := t.TempDir()
	fAdress := path.Join(tmp)
	store, err := newFileStore(fAdress, mockLogger)
	if err != nil {
		t.Errorf("New() = %v", err)
		return
	}
	sampleClient := wgmgmt.Client{
		PrivateKey:  "pri key",
		PublicKey:   "pub key",
		Name:        "ali",
		AllocatedIP: "192.168.10.2/32",
		AllowedIPs:  []string{"0.0.0.0/0"},
		QRCode:      "qcode",
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inErr := store.AddClient(sampleClient)
	if inErr != nil {
		t.Errorf("AddClient() = %v", inErr)
		return
	}
	sampleServer := wgmgmt.Server{
		Address:    "192.168.10.1",
		ListenPort: 514,
		PostUp:     []string{"post script-1", "post script-2"},
		PostDown:   []string{"pre script"},
		PrivateKey: "pri key",
		PublicKey:  "pub key",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	sErr := store.SaveServerConfig(sampleServer)
	if sErr != nil {
		t.Errorf("SaveServerConfig() = %v", sErr)
		return
	}
	tests := []struct {
		name    string
		cs      *fileStore
		wantErr bool
	}{
		{
			"save file",
			store,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cs.saveWgDataToFile(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigStore.saveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
