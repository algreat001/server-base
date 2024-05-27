package model

import (
	"encoding/json"
	"runtime/debug"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type SocketServer interface {
	On(event string, handler EventHandlerFn)
	OnSystem(event string, handler EventHandlerFn)
	Once(event string, handler EventHandlerFn)
	Off(event string, handler EventHandlerFn)
	SendError(event string, socket *SyncWebSocket, message string)
	SendMessagesToUser(userId *uuid.UUID, message *SocketResponseMessage)
	SendMessagesToUsers(userIds []*uuid.UUID, message *SocketResponseMessage)
	SendMessageToSocket(socket *SyncWebSocket, message *SocketResponseMessage)
	Broadcast(message *SocketResponseMessage, senderId *uuid.UUID, includeSender bool)
	StartSocketHandler(conn *websocket.Conn, user *User)
	AddSocketToRoom(roomId *uuid.UUID, socket *SyncWebSocket)
	RemoveSocketFromRoom(roomId *uuid.UUID, socket *SyncWebSocket)
	SendMessageToRoom(roomId *uuid.UUID, message *SocketResponseMessage)
}

type SocketService struct {
	authFn         AuthTestRouteFn
	userSockets    map[uuid.UUID][]*SyncWebSocket
	roomSockets    map[uuid.UUID][]*SyncWebSocket
	handlers       map[string][]*EventHandlerFn
	systemHandlers map[string][]*EventHandlerFn
	handlersOnce   map[string][]*EventHandlerFn
}

type SocketResponseMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type SocketResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type SocketReceiveMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func NewSocketService(authFn AuthTestRouteFn) SocketServer {
	return &SocketService{
		authFn:         authFn,
		userSockets:    make(map[uuid.UUID][]*SyncWebSocket),
		roomSockets:    make(map[uuid.UUID][]*SyncWebSocket),
		handlers:       make(map[string][]*EventHandlerFn),
		systemHandlers: make(map[string][]*EventHandlerFn),
		handlersOnce:   make(map[string][]*EventHandlerFn),
	}
}

func (s *SocketService) AddUserSocket(userUUID uuid.UUID, conn *websocket.Conn) *SyncWebSocket {
	if s.userSockets[userUUID] == nil {
		logrus.Info("Create new user socket list for userUUID: ", userUUID)
		s.userSockets[userUUID] = make([]*SyncWebSocket, 0)
	}

	socket := &SyncWebSocket{mutex: sync.Mutex{}, conn: conn}
	s.userSockets[userUUID] = append(s.userSockets[userUUID], socket)
	logrus.Info("Add new user socket for userUUID: ", s.userSockets[userUUID])
	return socket
}

func (s *SocketService) AddSocketToRoom(roomId *uuid.UUID, socket *SyncWebSocket) {
	if s.roomSockets[*roomId] == nil {
		logrus.Info("Create new room socket list for roomId: ", *roomId)
		s.roomSockets[*roomId] = make([]*SyncWebSocket, 0)
	}
	s.roomSockets[*roomId] = append(s.roomSockets[*roomId], socket)
	logrus.Info("Add new room socket for roomId: ", s.roomSockets[*roomId])
}
func (s *SocketService) RemoveSocketFromRoom(roomId *uuid.UUID, socket *SyncWebSocket) {
	for i, c := range s.roomSockets[*roomId] {
		if c == socket {
			s.roomSockets[*roomId] = append(s.roomSockets[*roomId][:i], s.roomSockets[*roomId][i+1:]...)
			break
		}
	}
	logrus.Info("Remove socket from room: ", *roomId)
	if len(s.roomSockets[*roomId]) == 0 {
		delete(s.roomSockets, *roomId)
		logrus.Info("Room is empty, remove room: ", *roomId)
	}
}

func (s *SocketService) removeSocketFromAllRooms(socket *SyncWebSocket) {
	for roomId := range s.roomSockets {
		s.RemoveSocketFromRoom(&roomId, socket)
	}
}

func (s *SocketService) SendMessageToRoom(roomId *uuid.UUID, message *SocketResponseMessage) {
	for _, socket := range s.roomSockets[*roomId] {
		s.SendMessageToSocket(socket, message)
	}
}

func (s *SocketService) CloseConnection(userUUID uuid.UUID, conn *websocket.Conn) {
	defer conn.Close()
	for i, c := range s.userSockets[userUUID] {
		if c.conn == conn {
			s.userSockets[userUUID] = append(s.userSockets[userUUID][:i], s.userSockets[userUUID][i+1:]...)
			s.removeSocketFromAllRooms(c)
			break
		}
	}
}

func (s *SocketService) GetConnSockets(userUUID *uuid.UUID) []*SyncWebSocket {
	return s.userSockets[*userUUID]
}

func (s *SocketService) StartSocketHandler(conn *websocket.Conn, user *User) {
	userUUID := user.Id
	ctx := NewSocketContext(user, s.AddUserSocket(*userUUID, conn))
	s.SendMessagesToUser(userUUID, &SocketResponseMessage{Event: "connected", Data: &SocketResponse{Status: "ok", Data: user}})

	go func() {
		defer s.CloseConnection(*userUUID, conn)
		s.runSystemHandler("on-start", ctx)
		for {
			_, messageRaw, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					logrus.Info("WebSocket connection closed:", err)
				} else {
					logrus.Error("Failed to read message:", err)
				}
				break
			}
			var message SocketReceiveMessage
			err = json.Unmarshal(messageRaw, &message)
			s.showMessage(message)

			if err != nil {
				logrus.Error("Failed to unmarshal socket message:", err)
				continue
			}

			s.runHandler(message.Event, &s.handlersOnce, ctx, message.Data)
			s.runHandler(message.Event, &s.handlers, ctx, message.Data)

			if s.handlersOnce[message.Event] != nil {
				for _, handler := range s.handlersOnce[message.Event] {
					s.removeHandler(message.Event, *handler, &s.handlersOnce)
				}
			}
		}
	}()
}

func (s *SocketService) showMessage(message SocketReceiveMessage) {
	data := make(map[string]any)
	err := json.Unmarshal(message.Data, &data)
	if err != nil {
		logrus.Error("Failed to unmarshal socket message data:", err)
		return
	}
	logrus.Info("Receive message - event:", message.Event, " data:", data)
}

func (s *SocketService) runSystemHandler(event string, ctx *SocketContext) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("Recovered in socket runSystemHandler; event: ", event, " recover: ", r, " stack: ", string(debug.Stack()))
		}
	}()
	if s.systemHandlers[event] != nil {
		for _, handler := range s.systemHandlers[event] {
			(*handler)(event, ctx, nil)
		}
	}
}

func (s *SocketService) runHandler(event string, handlersPointer *map[string][]*EventHandlerFn, ctx *SocketContext, args []byte) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("Recovered in socket runHandler; event: ", event, " recover: ", r, " stack: ", string(debug.Stack()))
		}
	}()
	user := ctx.GetUser()
	handlers := *handlersPointer
	if handlers[event] != nil {
		for _, handler := range handlers[event] {
			if !s.authFn(event, user) {
				logrus.Info("User ", user, " is not authorized for event ", event)
				continue
			}
			(*handler)(event, ctx, args)
		}
	}
}

func (s *SocketService) addHandler(event string, handler EventHandlerFn, handlersPointer *map[string][]*EventHandlerFn) {
	handlers := *handlersPointer
	if handlers[event] == nil {
		handlers[event] = make([]*EventHandlerFn, 0)
	}
	handlers[event] = append(handlers[event], &handler)
	logrus.Info("Add new handler for event '", event, "' : ", handler)
}

func (s *SocketService) removeHandler(event string, handler EventHandlerFn, handlersPointer *map[string][]*EventHandlerFn) {
	handlers := *handlersPointer
	for i, h := range handlers[event] {
		if h == &handler {
			handlers[event] = append(handlers[event][:i], handlers[event][i+1:]...)
			break
		}
	}
	logrus.Info("Remove handler for event: ", handler)
}

func (s *SocketService) On(event string, handler EventHandlerFn) {
	s.addHandler(event, handler, &s.handlers)
}
func (s *SocketService) OnSystem(event string, handler EventHandlerFn) {
	s.addHandler(event, handler, &s.systemHandlers)
}
func (s *SocketService) Once(event string, handler EventHandlerFn) {
	s.addHandler(event, handler, &s.handlersOnce)
}
func (s *SocketService) Off(event string, handler EventHandlerFn) {
	s.removeHandler(event, handler, &s.handlers)
	s.removeHandler(event, handler, &s.handlersOnce)
}

func (s *SocketService) SendMessagesToUser(userId *uuid.UUID, message *SocketResponseMessage) {
	for _, socket := range s.GetConnSockets(userId) {
		s.SendMessageToSocket(socket, message)
	}
}
func (s *SocketService) SendMessageToSocket(socket *SyncWebSocket, message *SocketResponseMessage) {
	if socket == nil || socket.conn == nil {
		logrus.Error("Socket is nil")
		return
	}
	socket.mutex.Lock()
	err := socket.conn.WriteJSON(message)
	socket.mutex.Unlock()
	if err != nil {
		logrus.Error("Failed to write message:", err)
	}
}

func (s *SocketService) SendError(event string, socket *SyncWebSocket, error string) {
	s.SendMessageToSocket(socket, &SocketResponseMessage{
		Event: event,
		Data:  map[string]any{"status": "error", "error": error},
	})
}

func (s *SocketService) SendMessagesToUsers(userIds []*uuid.UUID, message *SocketResponseMessage) {
	for _, userUUID := range userIds {
		s.SendMessagesToUser(userUUID, message)
	}
}

func (s *SocketService) Broadcast(message *SocketResponseMessage, senderId *uuid.UUID, includeSender bool) {
	for userUUID := range s.userSockets {
		if !includeSender && userUUID == *senderId {
			continue
		}
		s.SendMessagesToUser(&userUUID, message)
	}
}
