package main

import (
	"github.com/saiset-co/saiStorageMongo/src/sai/db/mongo"
	"github.com/saiset-co/saiStorageMongo/src/sai/network/http"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/api"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/settings"
)

func main() {
	settings.LoadSettings()
	settings.SaveSettings()
	mongo.SetMongoDBInv(settings.Settings.DB)

	if !settings.Settings.DB.Atlas.Enabled {
		mongo.StartMongod()
	}

	mongo.TestMongoConnection()
	http.SetHttpServerInv(settings.Settings.HttpServer)
	api.InitAPI()
	api.CreateDefaultRoles(settings.Settings.Auth.DefaultRoles)
	api.CreateDefaultSuperAdmin(settings.Settings.Auth.SuperAdmin)
	api.SyncTokens()
	http.SaiHttpServer.Start()
}
