package busmodule

import (
	"github.com/google/uuid"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type BusService struct {
	store   store.Store
	auth    *middleware.Auth
	busUser *model.User
}

func NewService(store store.Store, auth *middleware.Auth) *BusService {
	userId, _ := uuid.FromBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	user := &model.User{
		Id:    &userId,
		Email: "bus@rabbit.mq",
	}
	return &BusService{
		store:   store,
		auth:    auth,
		busUser: user,
	}
}

func (sc *BusService) addUserToService(userInfo *UserFromBus) error {
	user := &model.User{
		Id:    userInfo.UserId,
		Email: userInfo.Email,
	}
	err := sc.store.User().AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (sc *BusService) removeUserToService(userInfo *UserFromBus) error {
	user := &model.User{
		Id: userInfo.UserId,
	}
	err := sc.store.User().Remove(sc.busUser, user)
	if err != nil {
		return err
	}
	return nil
}
