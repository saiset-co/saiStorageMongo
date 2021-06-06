package main

import (
	"saiStorageMongo/src/sai/db/mongo"
	"saiStorageMongo/src/sai/network/http"
	"saiStorageMongo/src/sai_storage/api"
	"saiStorageMongo/src/sai_storage/settings"
)

func main() {
	settings.LoadSettings()
	settings.SaveSettings()
	mongo.SetMongoDBInv(settings.Settings.DB)
	mongo.StartMongod()
	mongo.TestMongodConnection()
	http.SetHttpServerInv(settings.Settings.HttpServer)
	api.InitAPI()
	api.CreateDefaultRoles(settings.Settings.Auth.DefaultRoles)
	api.CreateDefaultSuperAdmin(settings.Settings.Auth.SuperAdmin)
	api.SyncTokens()
	http.SaiHttpServer.Start()
}
