package auth

import "github.com/webmakom-com/mycointainer/src/Storage/src/sai/storage"

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
