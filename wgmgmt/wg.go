package wgmgmt

import "time"

// Client model
type Client struct {
	ID              string    `json:"id"`
	PrivateKey      string    `json:"private_key"`
	PublicKey       string    `json:"public_key"`
	Name            string    `json:"name" validate:"required"`
	AllocatedIP     string    `json:"allocated_ip" validate:"required,ip"`
	AllocatedIpCIDR string    `json:"allocated_ip_cidr"`
	AllowedIPs      []string  `json:"allowed_ips" validate:"required,iplist"`
	QRCode          string    `json:"qr_code"`
	DNSAddress      string    `json:"dns_address" validate:"required,ip"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ClientStat model
type ClientStat struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	IP           string    `json:"ip"`
	PrivateKey   string    `json:"private_key"`
	PublicKey    string    `json:"public_key"`
	HandShake    HandShake `json:"handshake"`
	AllocatedIPs string    `json:"allocated_ip"`
	SendByte     int64     `json:"send_byte"`
	ReceiveByte  int64     `json:"receive_byte"`
}

type HandShake struct {
	LastHandshakeTime time.Time `json:"last_handshake_time"`
	Seen              bool      `json:"seen"`
}

// Server model
type Server struct {
	Address            string    `json:"address" validate:"required,ip"`
	TunnelAddress      string    `json:"tunnel_address" validate:"required,cidr"`
	TunnelAddressMask  int       `json:"tunnel_address_mask"`
	TunnelAddressBits  int       `json:"tunnel_address_bits"`
	Interface          string    `json:"interface"`
	ListenPort         int       `json:"listen_port,string"` // ,string to get listen_port string input as int
	PostUp             []string  `json:"post_up"`
	PostDown           []string  `json:"post_down"`
	AutoGenerateKey    bool      `json:"auto_generate_key"`
	AutoGenerateScript bool      `json:"auto_generate_script"`
	ManualIP           bool      `json:"manual_ip"`
	PrivateKey         string    `json:"private_key"`
	PublicKey          string    `json:"public_key"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ServerStatus struct {
	InterfaceName string `json:"interface_name"`
	IsEnable      bool   `json:"is_enable"`
	AllClient     int    `json:"all_client"`
	EnableClient  int    `json:"enable_client"`
	Send          int64  `json:"send"`
	Receive       int64  `json:"receive"`
}
