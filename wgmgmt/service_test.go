package wgmgmt

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	type args struct {
		cfgFullPath   string
		configStorage IConfigStorage
	}
	tests := []struct {
		name string
		args args
		want IWgmgmtService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfgFullPath, "wg0", tt.args.configStorage, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wgService_GetAllClient(t *testing.T) {
	tests := []struct {
		name string
		wg   *wgService
		want []Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wg.GetAllClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wgService.GetAllClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wgService_GetClient(t *testing.T) {
	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		wg      *wgService
		args    args
		want    Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.wg.GetClient(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("wgService.GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wgService.GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wgService_AddClient(t *testing.T) {
	type args struct {
		newClient  Client
		autoKeyGen bool
	}
	tests := []struct {
		name    string
		wg      *wgService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.wg.AddClient(tt.args.newClient, tt.args.autoKeyGen); (err != nil) != tt.wantErr {
				t.Errorf("wgService.AddClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_wgService_UpdateClient(t *testing.T) {
	type args struct {
		updateClient Client
	}
	tests := []struct {
		name    string
		wg      *wgService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.wg.UpdateClient(tt.args.updateClient); (err != nil) != tt.wantErr {
				t.Errorf("wgService.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_wgService_RemoveClient(t *testing.T) {
	type args struct {
		rmClient Client
	}
	tests := []struct {
		name    string
		wg      *wgService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.wg.RemoveClient(tt.args.rmClient); (err != nil) != tt.wantErr {
				t.Errorf("wgService.RemoveClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_wgService_GetServer(t *testing.T) {
	tests := []struct {
		name string
		wg   *wgService
		want Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wg.GetServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wgService.GetServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wgService_UpdateServer(t *testing.T) {
	type args struct {
		serverfg Server
	}
	tests := []struct {
		name    string
		wg      *wgService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.wg.UpdateServer(tt.args.serverfg); (err != nil) != tt.wantErr {
				t.Errorf("wgService.UpdateServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_wgService_AddServer(t *testing.T) {
	//---- user mock for store ------
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockIStore := NewMockIConfigStorage(mockCtrl)
	mockIStore.EXPECT().SaveServerConfig(gomock.Any()).Return(nil).
		Do(func(x interface{}) {
			mockIStore.EXPECT().GetServerConfig().Return(x).Times(2)
			mockIStore.EXPECT().AllClient().Return(nil).Times(1)
		})
	//-----------------------------------------
	newWg := New("./wg-1.confg", "wg-1", mockIStore, nil)
	sampleServer := Server{
		Address:            "192.168.10.1",
		ListenPort:         514,
		AutoGenerateScript: true,
		AutoGenerateKey:    true,
		PostUp:             []string{"post script-1", "post script-2"},
		PostDown:           []string{"pre script"},
		PrivateKey:         "pri key",
		PublicKey:          "pub key",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	//-----------------------------------------
	type args struct {
		newServer  Server
		autoKeyGen bool
	}
	tests := []struct {
		name    string
		wg      *wgService
		args    args
		wantErr bool
	}{
		{
			"add server with ID",
			newWg,
			args{
				newServer:  sampleServer,
				autoKeyGen: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.wg.AddServer(tt.args.newServer, tt.args.autoKeyGen); (err != nil) != tt.wantErr {
				t.Errorf("wgService.AddServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
