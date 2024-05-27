package socketmodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/middleware"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
)

type SocketService struct {
	socketServer model.SocketServer
}

func NewService(auth *middleware.Auth) model.SocketServer {
	return model.NewSocketService(auth.TestAuthRoute)
}
