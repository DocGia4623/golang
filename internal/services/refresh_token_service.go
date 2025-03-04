package services

import (
	"fmt"
	"testwire/config"
	"testwire/internal/models"
	"testwire/internal/repository"
	"testwire/utils"
)

type RefreshTokenService interface {
	SaveToken(models.RefreshToken) error
	DeleteToken(token string) error
	FindToken(token string) (*models.RefreshToken, error)
	RefreshToken(token string, signedKey string) (string, string, error)
}

type RefreshTokenServiceImpl struct {
	RefreshTokenRepository repository.RefreshTokenRepository
}

func NewRefreshTokenServiceImpl(refreshTokenRepository repository.RefreshTokenRepository) RefreshTokenService {
	return &RefreshTokenServiceImpl{
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (a *RefreshTokenServiceImpl) RefreshToken(token string, signedKey string) (string, string, error) {
	// Kiểm tra refresh token có trong database không
	_, err := a.RefreshTokenRepository.FindByToken(token)
	if err != nil {
		return "", "", fmt.Errorf("cant find refresh token")
	}
	// Kiểm tra token có hợp lệ không
	config, _ := config.LoadConfig()
	sub, expRaw, err := utils.ValidateRefreshToken(token, config.RefreshTokenSecret)
	if err != nil {
		return "", "", fmt.Errorf("refresh token is invalid: %w", err)
	}

	// Xóa refresh token cũ khỏi database
	a.DeleteToken(token)

	// Chuyển `expRaw` từ `interface{}` về `int64`
	expFloat, ok := expRaw.(float64)
	if !ok {
		return "", "", fmt.Errorf("invalid token expiration format")
	}
	exp := int64(expFloat) // Chuyển thành kiểu int64 (UNIX timestamp)

	// Tạo access token mới
	accessToken, err := utils.GenerateAccessToken(config.AccessTokenExpiresIn, sub, config.AccessTokenSecret)
	if err != nil {
		return "", "", fmt.Errorf("cannot generate access token")
	}
	// Tạo refresh token mới
	newRefreshToken, err := utils.GenerateRefreshToken(exp, sub, config.RefreshTokenSecret)
	if err != nil {
		return "", "", fmt.Errorf("cannot generate refresh token")
	}

	// Lưu refresh token mới vào database
	a.SaveToken(models.RefreshToken{
		Token: newRefreshToken,
	})

	return accessToken, newRefreshToken, nil
}

func (a *RefreshTokenServiceImpl) SaveToken(refreshToken models.RefreshToken) error {
	return a.RefreshTokenRepository.Save(refreshToken)
}
func (a *RefreshTokenServiceImpl) FindToken(token string) (*models.RefreshToken, error) {
	return a.RefreshTokenRepository.FindByToken(token)
}
func (a *RefreshTokenServiceImpl) DeleteToken(token string) error {
	return a.RefreshTokenRepository.Delete(token)
}
