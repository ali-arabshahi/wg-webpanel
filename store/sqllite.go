package store

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"path"
	"strings"
	"time"
	"wireguard-web-ui/trafficcapture"
	"wireguard-web-ui/usermanagment"
	"wireguard-web-ui/wgmgmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	//"database/sql" package use this
	// _ "github.com/mattn/go-sqlite3"
)

// var (
// 	createTableQuery = []string{`CREATE TABLE IF NOT EXISTS "User" (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		name TEXT NULL,
// 		age int NULL,
// 		created DATE NULL
// 	)`}
// )

type sqliteStore struct {
	DBClient *gorm.DB
	dataPath string
	logger   ILogger
}

// func newSqliteStore(dataPath string, logger ILogger) (*sqliteStore, error) {
// 	dataAddr := path.Join(dataPath, "data.sql")
// 	dbStore := sqliteStore{
// 		dataPath: dataPath,
// 		logger:   logger,
// 	}
// 	var err error
// 	dbStore.DBClient, err = sql.Open("sqlite3", dataAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &dbStore, nil
// }

func newSqliteStore(dataPath string, logger ILogger) (*sqliteStore, error) {
	dataAddr := path.Join(dataPath, "data.sql")
	dbStore := sqliteStore{
		dataPath: dataPath,
		logger:   logger,
	}
	var err error
	dbStore.DBClient, err = gorm.Open(sqlite.Open(dataAddr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if !dbStore.DBClient.Migrator().HasTable(&usermanagment.AdminUser{}) {
		err := dbStore.DBClient.AutoMigrate(usermanagment.AdminUser{})
		if err != nil {
			return nil, err
		}
		clearPassword := "wireguardadmin"
		md5Hasher := md5.New()
		md5Hasher.Write([]byte(clearPassword))
		md5Password := hex.EncodeToString(md5Hasher.Sum(nil))
		admin := usermanagment.AdminUser{
			UserName: "wgadmin",
			Password: md5Password,
		}
		dbStore.DBClient.Create(&admin)
	}
	if !dbStore.DBClient.Migrator().HasTable(&trafficcapture.SourceGroup{}) {
		err := dbStore.DBClient.AutoMigrate(&trafficcapture.SourceGroup{})
		if err != nil {
			return nil, err
		}
	}
	if !dbStore.DBClient.Migrator().HasTable(&wgmgmt.Server{}) {
		err := dbStore.DBClient.AutoMigrate(&wgmgmt.Server{})
		if err != nil {
			return nil, err
		}
	}
	if !dbStore.DBClient.Migrator().HasTable(&wgmgmt.Client{}) {
		err := dbStore.DBClient.AutoMigrate(&wgmgmt.Client{})
		if err != nil {
			return nil, err
		}
	}
	return &dbStore, nil
}

//------------------------------------
func (sq *sqliteStore) CheckUserAndPassword(adminUser usermanagment.AdminUser) bool {
	user := usermanagment.AdminUser{}
	clearPassword := adminUser.Password
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(clearPassword))
	md5Password := hex.EncodeToString(md5Hasher.Sum(nil))
	result := sq.DBClient.Where(usermanagment.AdminUser{UserName: adminUser.UserName, Password: md5Password}).First(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return false
	}
	return true
}

//------------------------------------
func (sq *sqliteStore) GetServerConfig() wgmgmt.Server {
	result := struct {
		Address            string
		ID                 int64
		TunnelAddress      string
		TunnelAddressMask  int
		TunnelAddressBits  int
		Interface          string
		ListenPort         int
		PostUp             string
		PostDown           string
		AutoGenerateKey    bool
		AutoGenerateScript bool
		ManualIP           bool
		PrivateKey         string
		PublicKey          string
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}{}
	sq.DBClient.Model(&wgmgmt.Server{}).First(&result)
	server := wgmgmt.Server{
		Address:            result.Address,
		ID:                 result.ID,
		TunnelAddress:      result.TunnelAddress,
		TunnelAddressMask:  result.TunnelAddressMask,
		TunnelAddressBits:  result.TunnelAddressBits,
		Interface:          result.Interface,
		ListenPort:         result.ListenPort,
		PostUp:             strings.Split(result.PostUp, ","),
		PostDown:           strings.Split(result.PostDown, ","),
		AutoGenerateKey:    result.AutoGenerateKey,
		AutoGenerateScript: result.AutoGenerateScript,
		ManualIP:           result.ManualIP,
		PrivateKey:         result.PrivateKey,
		PublicKey:          result.PublicKey,
		CreatedAt:          result.CreatedAt,
		UpdatedAt:          result.UpdatedAt,
	}
	return server
}

//------------------------------------
func (sq *sqliteStore) SaveServerConfig(serverConfig wgmgmt.Server) error {
	server := sq.GetServerConfig()
	convertedConfig := struct {
		Address            string
		ID                 int64
		TunnelAddress      string
		TunnelAddressMask  int
		TunnelAddressBits  int
		Interface          string
		ListenPort         int
		PostUp             string
		PostDown           string
		AutoGenerateKey    bool
		AutoGenerateScript bool
		ManualIP           bool
		PrivateKey         string
		PublicKey          string
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}{
		serverConfig.Address,
		serverConfig.ID,
		serverConfig.TunnelAddress,
		serverConfig.TunnelAddressMask,
		serverConfig.TunnelAddressBits,
		serverConfig.Interface,
		serverConfig.ListenPort,
		strings.Join(serverConfig.PostUp, ","),
		strings.Join(serverConfig.PostDown, ","),
		serverConfig.AutoGenerateKey,
		serverConfig.AutoGenerateScript,
		serverConfig.ManualIP,
		serverConfig.PrivateKey,
		serverConfig.PublicKey,
		serverConfig.CreatedAt,
		serverConfig.UpdatedAt,
	}
	convertedConfig.ID = server.ID
	res := sq.DBClient.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Model(&wgmgmt.Server{}).Create(&convertedConfig)
	// var count int64
	// sq.DBClient.Model(wgmgmt.Server{}).Count(&count)
	// var dErr error
	// server := sq.GetServerConfig()
	// if server.ID == 0 {
	// 	res := sq.DBClient.Model(wgmgmt.Server{}).Create(&convertedConfig)
	// 	if res.Error != nil {
	// 		dErr = res.Error
	// 	}
	// } else {
	// 	// convertedConfig.ID = server.ID
	// 	res := sq.DBClient.Model(wgmgmt.Server{}).Where("id = ?", server.ID).Save(&convertedConfig)
	// 	if res.Error != nil {
	// 		dErr = res.Error
	// 	}
	// }
	return res.Error
}

//------------------------------------
func (sq *sqliteStore) AllClient() []wgmgmt.Client {
	users := []wgmgmt.Client{}
	type client struct {
		ID              int64
		PrivateKey      string
		PublicKey       string
		Name            string
		AllocatedIP     string
		AllocatedIpCIDR string
		AllowedIPs      string
		QRCode          string
		DNSAddress      string
		Enabled         bool
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}
	allclient := []client{}
	result := sq.DBClient.Model(users).Find(&allclient)
	if result.RowsAffected > 0 {
		for _, item := range allclient {
			users = append(users, wgmgmt.Client{
				ID:              item.ID,
				PrivateKey:      item.PrivateKey,
				PublicKey:       item.PublicKey,
				Name:            item.Name,
				AllocatedIP:     item.AllocatedIP,
				AllocatedIpCIDR: item.AllocatedIpCIDR,
				AllowedIPs:      strings.Split(item.AllowedIPs, ","),
				QRCode:          item.QRCode,
				DNSAddress:      item.DNSAddress,
				Enabled:         item.Enabled,
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
			})
		}
	}
	return users
}

//------------------------------------
func (sq *sqliteStore) ClientByID(clientID int64) (wgmgmt.Client, error) {
	cl := struct {
		ID              int64
		PrivateKey      string
		PublicKey       string
		Name            string
		AllocatedIP     string
		AllocatedIpCIDR string
		AllowedIPs      string
		QRCode          string
		DNSAddress      string
		Enabled         bool
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}{}
	result := sq.DBClient.Model(wgmgmt.Client{}).First(&cl, clientID)
	if result.Error != nil {
		return wgmgmt.Client{}, result.Error
	}
	client := wgmgmt.Client{
		ID:              cl.ID,
		PrivateKey:      cl.PrivateKey,
		PublicKey:       cl.PublicKey,
		Name:            cl.Name,
		AllocatedIP:     cl.AllocatedIP,
		AllocatedIpCIDR: cl.AllocatedIpCIDR,
		AllowedIPs:      strings.Split(cl.AllowedIPs, ","),
		QRCode:          cl.QRCode,
		DNSAddress:      cl.DNSAddress,
		Enabled:         cl.Enabled,
		CreatedAt:       cl.CreatedAt,
		UpdatedAt:       cl.UpdatedAt,
	}
	return client, nil
}

//------------------------------------
func (sq *sqliteStore) ClientTrafficByIP(clientIP string) (SendByte int64, ReceiveByte int64, Error error) {
	clientTraffic := trafficcapture.SourceGroup{}
	result := sq.DBClient.Where(trafficcapture.SourceGroup{SourceIP: clientIP}).First(&clientTraffic)
	if result.Error != nil {
		return 0, 0, result.Error
	}
	return clientTraffic.SendByte, clientTraffic.RecieveByte, nil
}

//------------------------------------
func (sq *sqliteStore) ClientNumber() (int, int) {
	var enabledUser int64
	var allUser int64
	sq.DBClient.Model(&wgmgmt.Client{}).Where("enabled = ?", true).Count(&enabledUser)
	sq.DBClient.Model(&wgmgmt.Client{}).Count(&allUser)
	return int(allUser), int(enabledUser)
}

//------------------------------------
func (sq *sqliteStore) AllTraffic() (SendByte int64, ReceiveByte int64, Error error) {
	type traffic struct {
		Send    int64
		Receive int64
	}
	allTraffic := traffic{}
	qErr := sq.DBClient.Model(&trafficcapture.SourceGroup{}).Select("sum(send_byte) as Send,sum(recieve_byte) as receive").First(&allTraffic)
	if qErr.Error != nil {
		return 0, 0, qErr.Error
	}
	return allTraffic.Send, allTraffic.Receive, nil
}

//------------------------------------
func (sq *sqliteStore) AddClient(newClient wgmgmt.Client) error {
	cl := struct {
		ID              int64
		PrivateKey      string
		PublicKey       string
		Name            string
		AllocatedIP     string
		AllocatedIpCIDR string
		AllowedIPs      string
		QRCode          string
		DNSAddress      string
		Enabled         bool
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}{
		newClient.ID,
		newClient.PrivateKey,
		newClient.PublicKey,
		newClient.Name,
		newClient.AllocatedIP,
		newClient.AllocatedIpCIDR,
		strings.Join(newClient.AllowedIPs, ","),
		newClient.QRCode,
		newClient.DNSAddress,
		newClient.Enabled,
		newClient.CreatedAt,
		newClient.UpdatedAt,
	}

	result := sq.DBClient.Model(wgmgmt.Client{}).Create(&cl)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//------------------------------------
func (sq *sqliteStore) UsedPrivateIPS() []string {
	usedIp := []string{}
	// get all client ip
	rows, err := sq.DBClient.Model(&wgmgmt.Client{}).Select("allocated_ip").Group("allocated_ip").Rows()
	if err != nil {
		fmt.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ip string
			rows.Scan(&ip)
			usedIp = append(usedIp, ip)
		}
	}
	//----------------------
	server := sq.GetServerConfig()
	// server := wgmgmt.Server{}
	// sErr := sq.DBClient.Find(&server)
	// if sErr.Error != nil {
	// 	fmt.Println(sErr.Error)
	if server.ID == 0 {
		fmt.Println("no server config")
	} else {
		serverIPRange := server.TunnelAddress
		ipSrv, ipnet, pErr := net.ParseCIDR(serverIPRange)
		if pErr != nil {
			fmt.Println(pErr.Error())
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
	}
	return usedIp
}

//------------------------------------
func (sq *sqliteStore) UpdateClient(client wgmgmt.Client) error {
	cl := struct {
		ID              int64
		PrivateKey      string
		PublicKey       string
		Name            string
		AllocatedIP     string
		AllocatedIpCIDR string
		AllowedIPs      string
		QRCode          string
		DNSAddress      string
		Enabled         bool
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}{
		client.ID,
		client.PrivateKey,
		client.PublicKey,
		client.Name,
		client.AllocatedIP,
		client.AllocatedIpCIDR,
		strings.Join(client.AllowedIPs, ","),
		client.QRCode,
		client.DNSAddress,
		client.Enabled,
		client.CreatedAt,
		client.UpdatedAt,
	}
	result := sq.DBClient.Model(wgmgmt.Client{}).Select("*").Where("id = ?", cl.ID).Updates(&cl)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//------------------------------------
func (sq *sqliteStore) RemoveClient(client wgmgmt.Client) error {
	result := sq.DBClient.Delete(&client, client.ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//------------------------------------
func (sq *sqliteStore) WriteSourceGroup(sourceIps map[string]*trafficcapture.SourceGroup) {
	for ip, value := range sourceIps {
		// res := db.Model(&user).Where("name = ?", user.Name).Updates(map[string]interface{}{"age": gorm.Expr("age + ?", user.Age)})
		// PacketCount
		// SendByte
		// RecieveByte
		upRes := sq.DBClient.Model(&value).Where("source_ip = ?", ip).
			Updates(map[string]interface{}{
				"packet_count": gorm.Expr("packet_count + ?", value.PacketCount),
				"send_byte":    gorm.Expr("send_byte + ?", value.SendByte),
				"recieve_byte": gorm.Expr("recieve_byte + ?", value.RecieveByte),
			})
		if upRes.RowsAffected == 0 {
			inRes := sq.DBClient.Create(&value)
			if inRes.Error != nil {
				fmt.Println(inRes.Error)
			}
		}
	}
}

//**************************************************************//
//                      Private Function                        //
//**************************************************************//

// func (sq *sqliteStore) createTables(tablestruct []string) error {
// 	for _, table := range tablestruct {
// 		_, err := sq.DBClient.Exec(table)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
