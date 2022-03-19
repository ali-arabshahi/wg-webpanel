package usermanagment

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	TokenUserFieldName = "user"
	JwtExpiryFieldName = "exp"
)

type IUserManagment interface {
	GenerateToken(user string) (token string, err error)
	ValidateToken(token string) (tokenInfo map[string]interface{}, err error)
	CheckUserInfo(adminUser AdminUser) error
}

type IUserStore interface {
	CheckUserAndPassword(adminUser AdminUser) bool
}

type userMgmtService struct {
	tokenSecret string
	dateLayout  string
	userStore   IUserStore
}

func New(secret string, userStore IUserStore) *userMgmtService {
	return &userMgmtService{
		tokenSecret: secret,
		dateLayout:  "2006-01-02T15:04:05",
		userStore:   userStore,
	}
}

func (mgmt *userMgmtService) GenerateToken(user string) (string, error) {
	now := time.Now()
	exTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	timeFormatExp := exTime.Format(mgmt.dateLayout)
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		TokenUserFieldName: user,
		JwtExpiryFieldName: timeFormatExp,
	})

	token, err := tokenObj.SignedString([]byte(mgmt.tokenSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

//----------------------------------------------------------
func (mgmt *userMgmtService) ValidateToken(token string) (map[string]interface{}, error) {
	tokenObj, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected sign method: %v", t.Header["alg"])
			}
			return []byte(mgmt.tokenSecret), nil
		})
	if err != nil {
		return nil, err
	}
	tokeninfo := tokenObj.Claims.(jwt.MapClaims)
	isExpire := mgmt.isTokenExpire(tokeninfo)
	if isExpire {
		var emptyVal map[string]interface{}
		return emptyVal, fmt.Errorf("token is expire")
	}
	return tokeninfo, nil
}

//----------------------------------------------------------
func (mgmt *userMgmtService) CheckUserInfo(adminUser AdminUser) error {
	isValid := mgmt.userStore.CheckUserAndPassword(adminUser)
	if !isValid {
		return fmt.Errorf("user name or password is invali")
	}
	return nil
}

//-------------------------------------------------
//                Private functions
//-------------------------------------------------

// isTokenExpire : check if token exired
func (mgmt *userMgmtService) isTokenExpire(tokenInfo map[string]interface{}) bool {
	isExpire := false
	expireTimeStr := tokenInfo[JwtExpiryFieldName].(string)
	expireTime, err := time.ParseInLocation(mgmt.dateLayout, expireTimeStr, time.Local)
	if err != nil {
		return false
	}
	if expireTime.After(time.Now()) {
		isExpire = true
	}
	return isExpire
}
