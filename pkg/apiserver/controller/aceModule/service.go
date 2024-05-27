package acemodule

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/model"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

type AceService struct {
	store    store.Store
	executor *model.User
}

func NewService(store store.Store) *AceService {
	return &AceService{
		store: store,
	}
}

func (s *AceService) getACEGroups(_ *model.SocketContext) ([]*model.AceGroup, error) {
	return s.store.AccessControlEntity().GetGroups()
}

func (s *AceService) addACEGroup(ctx *model.SocketContext, groupReq *dto.ACEGroupReq) error {
	group := &model.AceGroup{
		Name: groupReq.Name,
	}

	_, err := s.store.AccessControlEntity().SetGroup(ctx.GetUser(), group)

	return err
}

func (s *AceService) removeACEGroup(ctx *model.SocketContext, groupReq *dto.ACEGroupReq) error {
	group := &model.AceGroup{
		Id: groupReq.Id,
	}

	return s.store.AccessControlEntity().RemoveGroup(ctx.GetUser(), group)
}
func (s *AceService) updateACEGroup(ctx *model.SocketContext, groupReq *dto.ACEGroupReq) error {
	group := &model.AceGroup{
		Id:   groupReq.Id,
		Name: groupReq.Name,
	}

	return s.store.AccessControlEntity().UpdateGroup(ctx.GetUser(), group)
}

func (s *AceService) getUserACEGroups(_ *model.SocketContext, userReq *dto.UserReq) ([]*model.AceGroup, error) {
	user := model.NewUserFromId(userReq.Id)
	groups, err := s.store.AccessControlEntity().GetUserGroups(user)
	if err != nil {
		return []*model.AceGroup{}, err
	}
	return groups, nil
}

func (s *AceService) addUserToACEGroup(ctx *model.SocketContext, userInGroupReq *dto.UserInAccessGroupReq) ([]*model.AceGroup, error) {
	user := model.NewUserFromId(userInGroupReq.UserId)
	group, err := s.store.AccessControlEntity().GetGroup(userInGroupReq.GroupId)
	if err != nil {
		return []*model.AceGroup{}, err
	}

	if err = s.store.AccessControlEntity().SetUserGroup(ctx.GetUser(), user, group); err != nil {
		return []*model.AceGroup{}, err
	}
	return s.store.AccessControlEntity().GetUserGroups(user)
}

func (s *AceService) removeUserFromACEGroup(ctx *model.SocketContext, userInGroupReq *dto.UserInAccessGroupReq) ([]*model.AceGroup, error) {
	user := model.NewUserFromId(userInGroupReq.UserId)
	group, err := s.store.AccessControlEntity().GetGroup(userInGroupReq.GroupId)
	if err != nil {
		return []*model.AceGroup{}, err
	}

	if err = s.store.AccessControlEntity().RemoveUserFromGroup(ctx.GetUser(), user, group); err != nil {
		return []*model.AceGroup{}, err
	}
	return s.store.AccessControlEntity().GetUserGroups(user)
}

func (s *AceService) getACEsForGroup(_ *model.SocketContext, userReq *dto.ACEGroupIdReq) ([]*model.Ace, error) {
	group, err := s.store.AccessControlEntity().GetGroup(userReq.Id)
	if err != nil {
		return []*model.Ace{}, err
	}

	return s.store.AccessControlEntity().GetGroupAces(group)
}

func (s *AceService) getACL(_ *model.SocketContext) ([]*model.Ace, error) {
	return s.store.AccessControlEntity().GetACL()
}

func (s *AceService) addACEToGroup(ctx *model.SocketContext, aceToGroupReq *dto.ACEToGroupReq) error {
	group, err := s.store.AccessControlEntity().GetGroup(aceToGroupReq.Group.Id)
	if err != nil {
		return err
	}
	ace := &model.Ace{Path: aceToGroupReq.Path}

	return s.store.AccessControlEntity().SetGroupAce(ctx.GetUser(), group, ace)
}

func (s *AceService) updateACE(ctx *model.SocketContext, updateACE *dto.ACEToGroupReq) error {

	ace := &model.Ace{Id: updateACE.Id, Path: updateACE.Path, Group: &model.AceGroup{Id: updateACE.Group.Id, Name: updateACE.Group.Name}}

	return s.store.AccessControlEntity().UpdateAce(ctx.GetUser(), ace)
}

func (s *AceService) removeACE(ctx *model.SocketContext, aceReq *dto.ACEIdReq) error {
	ace := &model.Ace{Id: aceReq.Id}
	return s.store.AccessControlEntity().RemoveAce(ctx.GetUser(), ace)
}
