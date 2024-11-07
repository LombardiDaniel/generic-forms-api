package services

type AuthService interface {
	Authenticate(key string) error
}
