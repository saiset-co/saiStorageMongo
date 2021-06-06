package auth

import "saiStorageMongo/src/sai/storage"

type Session struct {
	ID               string
	User             User
	Token            Token
	RoutePermissions Permission
}

func CreateSession() *Session {
	session := &Session{
		ID: storage.CreateDocumentID(),
	}
	return session
}
