package sqlstore

import (
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/tools"

	_ "github.com/lib/pq"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/store"
)

// Store for connect database
type Store struct {
	db                            tools.DataBase
	langRepository                *LangRepository
	accessControlEntityRepository *AccessControlEntityRepository
	logRepository                 *LogRepository
	userRepository                *UserRepository
	tokenRepository               *TokenRepository
}

// New database store
func New(db tools.DataBase) *Store {
	return &Store{
		db: db,
	}
}

type Repositories interface {
	LangRepository | AccessControlEntityRepository | LogRepository | UserRepository | TokenRepository
}

func getRepository[R Repositories](repo *R, s *Store) *R {
	if repo != nil {
		return repo
	}
	repo = &R{
		store: s,
	}
	return repo
}

func (s *Store) AccessControlEntity() store.AccessControlEntityRepository {
	return getRepository(s.accessControlEntityRepository, s)
}
func (s *Store) Lang() store.LangRepository {
	return getRepository(s.langRepository, s)
}
func (s *Store) Log() store.LogRepository {
	return getRepository(s.logRepository, s)
}
func (s *Store) User() store.UserRepository {
	return getRepository(s.userRepository, s)
}
func (s *Store) Token() store.TokenRepository {
	return getRepository(s.tokenRepository, s)
}
