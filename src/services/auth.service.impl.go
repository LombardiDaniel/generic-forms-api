package services

import "github.com/LombardiDaniel/generic-data-collector-api/utils"

type AuthServiceImpl struct {
	authStringsMap map[string]bool
}

func NewAuthServiceImpl(authStrings []string) AuthService {
	authStringMap := make(map[string]bool)
	for _, v := range authStrings {
		authStringMap[v] = true
	}

	return &AuthServiceImpl{
		authStringsMap: authStringMap,
	}
}

func (s *AuthServiceImpl) Authenticate(key string) error {
	_, exists := s.authStringsMap[key]

	if !exists {
		return utils.ErrAuth
	}

	return nil
}
