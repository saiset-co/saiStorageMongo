package auth

import (
	"github.com/saiset-co/saiStorageMongo/src/sai/storage"
)

var (
	RoleCollection = "roles"
)

type Role struct {
	ID          string                `json:"_id" bson:"-"`
	Name        string                `json:"name" bson:"name"`
	Permissions map[string]Permission `json:"permissions" bson:"permissions"`
}

func CreateNewRole() *Role {
	return &Role{
		ID: storage.CreateDocumentID(),
	}
}

func (role *Role) CreateID() {
	role.ID = storage.CreateDocumentID()
}

func (permission Rights) Sum() int {
	return permission.Read + permission.Write
}
