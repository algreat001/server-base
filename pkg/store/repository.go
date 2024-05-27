package store

import (
	"github.com/google/uuid"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
)

type AccessControlEntityRepository interface {
	//acess groups CRUD operation
	GetGroups() ([]*model.AceGroup, error)
	GetGroup(groupId *uuid.UUID) (*model.AceGroup, error)
	SetGroup(*model.User, *model.AceGroup) (*model.AceGroup, error)
	RemoveGroup(*model.User, *model.AceGroup) error
	UpdateGroup(*model.User, *model.AceGroup) error

	//user in acess groups CRUD operation
	GetUserGroups(*model.User) ([]*model.AceGroup, error)
	SetUserGroup(*model.User, *model.User, *model.AceGroup) error
	UpdateUserGroups(*model.User, *model.User) error
	RemoveUserFromGroup(*model.User, *model.User, *model.AceGroup) error

	//acess control entity CRUD operation
	GetGroupAces(*model.AceGroup) ([]*model.Ace, error)
	SetGroupAce(*model.User, *model.AceGroup, *model.Ace) error

	GetACL() ([]*model.Ace, error)
	RemoveAce(*model.User, *model.Ace) error
	UpdateAce(*model.User, *model.Ace) error
	SetAce(*model.User, *model.Ace) (*model.Ace, error)

	//operations for the user in relation to the acess control entity
	GetUserRightForPath(user *model.User, path string) (bool, error)
}

type LangRepository interface {
	GetLangList() ([]*model.Lang, error)
	GetLang(*model.Lang) (*model.Lang, error)
}

type LogRepository interface {
	GetLogList(executor *model.User, numPage int, sizePage int, order string, descending bool, filter string) ([]*dto.LogReq, int, error)
}

type UserRepository interface {
	AddUser(*model.User) error
	Update(*model.User) error
	Remove(*model.User, *model.User) error
	GetNumberAllUsers() (int, error)
	Find(*uuid.UUID) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	GetUsers(numPage int, sizePage int, order string, descending bool, filter string) ([]*model.User, error)
}

type TokenRepository interface {
	VerifyToken(token string) error
	Create(token *model.Token) error
	Delete(token string) error
}
