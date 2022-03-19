package wgmgmt

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_wgService_saveConfigToFile(t *testing.T) {
	//---- user mock for store ------
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockIStore := NewMockIConfigStorage(mockCtrl)
	//-------- sample server and client config ------------//
	sampleServer := Server{
		Address:    "192.168.10.1",
		ListenPort: 514,
		PostUp:     []string{"post script-1", "post script-2"},
		PostDown:   []string{"pre script"},
		PrivateKey: "pri key",
		PublicKey:  "pub key",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	sampleClient := []Client{
		Client{
			ID:          "id-client-1",
			PrivateKey:  "pri key",
			PublicKey:   "pub key",
			Name:        "ali",
			AllocatedIP: "192.168.10.2/32",
			AllowedIPs:  []string{"0.0.0.0/0"},
			QRCode:      "qcode",
			Enabled:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Client{
			ID:          "id-client-2",
			PrivateKey:  "pri ke2",
			PublicKey:   "pub key-2",
			Name:        "hasan",
			AllocatedIP: "192.168.10.3/32",
			AllowedIPs:  []string{"0.0.0.0/0"},
			QRCode:      "qcode",
			Enabled:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	//-----------------------------------------------------//
	mockIStore.EXPECT().AllClient().Return(sampleClient).Times(1)
	mockIStore.EXPECT().GetServerConfig().Return(sampleServer).Times(1)
	//-----------------------------------------------------//
	newWg := New("./wg-1.confg", "wg-1", mockIStore, nil)

	tests := []struct {
		name    string
		wg      *wgService
		wantErr bool
	}{
		{
			"write config",
			newWg,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.wg.saveConfigToFile(); (err != nil) != tt.wantErr {
				t.Errorf("wgService.saveConfigToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
