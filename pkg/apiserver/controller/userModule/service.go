package usermodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type UserService struct {
	store    store.Store
	executor *model.User
}

func NewService(store store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) addUser(ctx *model.SocketContext, reqUser *dto.UserUpdateReq) (*model.User, error) {
	user := model.NewUserFromDto(reqUser)
	if err := s.store.User().AddUser(user); err != nil {
		return nil, err
	}
	if err := s.store.AccessControlEntity().UpdateUserGroups(ctx.GetUser(), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) updateUser(ctx *model.SocketContext, reqUser *dto.UserUpdateReq) (*model.User, error) {
	findUser, err := s.store.User().FindByEmail(reqUser.Email)
	if err != nil || findUser == nil {
		return nil, err
	}

	user := model.NewUserFromDto(reqUser)
	if err := s.store.User().Update(user); err != nil {
		return nil, err
	}
	user.Id = findUser.Id
	if err := s.store.AccessControlEntity().UpdateUserGroups(ctx.GetUser(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) removeUser(ctx *model.SocketContext, user *dto.UserUpdateReq) error {
	u := &model.User{
		Id: user.Id,
	}
	if err := s.store.User().Remove(ctx.GetUser(), u); err != nil {
		return err
	}
	return nil
}

func (s *UserService) getUsers(_ *model.SocketContext, requestUser *dto.UsersPageReq) (*dto.UserPage, error) {
	users, err := s.store.User().GetUsers(requestUser.Page, requestUser.PageSize, requestUser.Order, *requestUser.Descending, requestUser.Filter)

	if err != nil {
		return nil, servererrors.ErrorInternal
	}

	count, err := s.store.User().GetNumberAllUsers()

	if err != nil {
		return nil, servererrors.ErrorInternal
	}

	for _, user := range users {
		groups, err := s.store.AccessControlEntity().GetUserGroups(user)
		if err != nil {
			return nil, servererrors.ErrorInternal
		}
		user.SetGroups(groups)
	}

	dtoUsers := make([]*dto.UserUpdateReq, len(users))
	for i, user := range users {
		dtoUsers[i] = user.ToDto()
	}

	return &dto.UserPage{
		CountAllRecord: count,
		Data:           dtoUsers,
	}, nil
}

func (s *UserService) getGroups(_ *model.SocketContext, user *model.User) ([]*model.AceGroup, error) {
	return s.store.AccessControlEntity().GetUserGroups(user)
}
