package busmodule

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type AdditionalCommandFn func(event string, user *model.User, params any)

type BusController struct {
	store              store.Store
	bus                *model.ServiceBus
	Service            *BusService
	additionalCommands map[string][]*AdditionalCommandFn
}

type UserFromBus struct {
	UserId            *uuid.UUID     `json:"user_id"`
	Email             string         `json:"user_email,omitempty"`
	ServiceResourceId *uuid.UUID     `json:"service_resource_id"`
	Params            map[string]any `json:"params,omitempty"`
}

func InitBusController(store store.Store, bus *model.ServiceBus, auth *middleware.Auth) *BusController {
	sc := &BusController{
		store:              store,
		Service:            NewService(store, auth),
		bus:                bus,
		additionalCommands: make(map[string][]*AdditionalCommandFn),
	}

	bus.On("add_user_to_service", sc.addUserToService)
	bus.On("remove_user_from_service", sc.removeUserFromService)

	return sc
}

func (sc *BusController) AddCommand(command string, handler AdditionalCommandFn) {
	if sc.additionalCommands[command] == nil {
		sc.additionalCommands[command] = make([]*AdditionalCommandFn, 0)
	}
	sc.additionalCommands[command] = append(sc.additionalCommands[command], &handler)
}

func (sc *BusController) executeCommand(event string, user *model.User, params any) {
	if handlers, ok := sc.additionalCommands[event]; ok {
		for _, handler := range handlers {
			(*handler)(event, user, params)
		}
	}
}

func (sc *BusController) addUserToService(event string, args []byte) {
	userInfo := &UserFromBus{}
	err := json.Unmarshal(args, userInfo)
	if err != nil {
		logrus.Info("Error add user to service (unmarshal) - ", err.Error())
		return
	}
	err = sc.Service.addUserToService(userInfo)
	if err != nil {
		logrus.Info("Error add user to service (store) - ", err.Error())

	}
	userInfo.Params["user_id"] = userInfo.UserId.String()
	sc.executeCommand(event, sc.Service.busUser, userInfo.Params)
}

func (sc *BusController) removeUserFromService(event string, args []byte) {
	userInfo := &UserFromBus{}
	err := json.Unmarshal(args, userInfo)
	logrus.Info("Remove user from service - ", userInfo)
	if err != nil {
		logrus.Info("Error remove user to service (unmarshal) - ", err.Error())
		return
	}
	err = sc.Service.removeUserToService(userInfo)
	if err != nil {
		logrus.Info("Error remove user to service (store) - ", err.Error())
	}
	userInfo.Params["user_id"] = userInfo.UserId.String()
	sc.executeCommand(event, sc.Service.busUser, userInfo.Params)
}
