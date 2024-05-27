package store

type Store interface {
	AccessControlEntity() AccessControlEntityRepository
	Lang() LangRepository
	Log() LogRepository
	User() UserRepository
	Token() TokenRepository
}
