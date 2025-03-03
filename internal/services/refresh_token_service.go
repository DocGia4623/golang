package services

import "testwire/internal/models"

type RefreshTokenService interface {
	SaveToken(models.RefreshToken) error
	DeleteToken(token string) error
	FindToken(token string) (*models.RefreshToken, error)
	RefreshToken(token string, signedKey string) (string, string, error)
}
