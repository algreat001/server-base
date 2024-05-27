package model

import (
	"context"
	"github.com/google/uuid"
	"sync"

	"github.com/gorilla/websocket"
)

type SyncWebSocket struct {
	mutex sync.Mutex
	conn  *websocket.Conn
}

type EventHandlerFn func(event string, ctx *SocketContext, args []byte)

type AuthTestRouteFn func(path string, user *User) bool

type SocketContext struct {
	user   *User
	ctx    context.Context
	socket *SyncWebSocket
	id     uuid.UUID
}

func NewSocketContext(user *User, socket *SyncWebSocket) *SocketContext {
	return &SocketContext{
		user:   user,
		socket: socket,
		ctx:    context.Background(),
		id:     uuid.New(),
	}
}

func (s *SocketContext) GetId() uuid.UUID {
	return s.id
}

func (s *SocketContext) GetUser() *User {
	return s.user
}

func (s *SocketContext) GetSocket() *SyncWebSocket {
	return s.socket
}

func (s *SocketContext) GetValue(key string) any {
	return s.ctx.Value(key)
}

func (s *SocketContext) SetValue(key string, value any) {
	s.ctx = context.WithValue(s.ctx, key, value)
}
