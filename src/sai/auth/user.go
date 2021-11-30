package auth

import (
	"crypto"
	"encoding/base64"
	"github.com/saiset-co/saiStorageMongo/src/sai/common"
	"github.com/saiset-co/saiStorageMongo/src/sai/storage"
)

var UserCollection = "users"

type User struct {
	ID       string      `json:"_id" bson:"-"`
	Email    string      `json:"email" bson:"email"`
	IPAddr   string      `json:"ipaddr" bson:"ipaddr"`
	Password string      `json:"password" bson:"password"`
	Data     interface{} `json:"data" bson:"data"`
	RoleID   string      `json:"role_id" bson:"role_id"`
}

func (user *User) SetCredentials(email string, password string) {
	user.Email = email
	user.Password = HashUserPassword(password)
}

func CreateNewUser() *User {
	return &User{
		ID: storage.CreateDocumentID(),
	}
}

func HashUserPassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(crypto.SHA256.New().Sum(crypto.MD5.New().Sum([]byte(password)))))
}

func (user *User) String() string {
	return string(common.ConvertInterfaceToJson(user))
}

func (user *User) Can(permission Permission) bool {
	return true
	//return user.Role.Permission.Sum() >= permission.Sum()
}
