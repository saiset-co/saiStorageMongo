package settings

import (
	"fmt"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/tkanos/gonfig"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/auth"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/network/http"
)

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type LocalMongoConfig struct {
	Config  DatabaseConfig `json:"config"`
	Host    string         `json:"host"`
	Enabled bool           `json:"enabled"`
}

type AtlasMongoConfig struct {
	Config           DatabaseConfig `json:"config"`
	Host             string         `json:"host"`
	Enabled          bool           `json:"enabled"`
	ConnectionString string         `json:"connection_string"`
}

type DBConfig struct {
	Local LocalMongoConfig `json:"local"`
	Atlas AtlasMongoConfig `json:"atlas"`
}

type SuperAdminConfig struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ParamConfig struct {
	Name  string   `json:"name"`
	Rules []string `json:"rules"`
	//Required bool     `json:"required"`
}

type RightsConfig struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type PermissionConfig struct {
	URL    string        `json:"url"`
	Rights auth.Rights   `json:"rights"`
	Params []ParamConfig `json:"params"`
}

type RoleConfig struct {
	Name        string             `json:"name" bson:"name"`
	Permissions []PermissionConfig `json:"permissions" bson:"permissions"`
}

type AuthConfig struct {
	Enabled      bool             `json:"enabled"`
	DefaultRoles []RoleConfig     `json:"default_roles"`
	SuperAdmin   SuperAdminConfig `json:"super_admin"`
}

type Configuration struct {
	HttpServer http.HttpServer `json:"http_server"`
	DB         DBConfig        `json:"db"`
	Auth       AuthConfig      `json:"auth"`
}

var (
	Settings = Configuration{}
)

func LoadSettings() {
	err := gonfig.GetConf("config.json", &Settings)
	if err != nil {
		fmt.Println(err)
	}
}

func SaveSettings() {

}
