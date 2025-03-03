package repository

import "testwire/internal/models"

type RefreshTokenRepository interface {
	Save(refreshToken models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	Delete(token string) error
}
