package repository

import (
	"testwire/internal/models"

	"gorm.io/gorm"
)

type refreshTokenRepositoryImpl struct {
	Db *gorm.DB
}

func NewRefreshTokenRepositoryImpl(Db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{Db: Db}
}

func (r *refreshTokenRepositoryImpl) Save(refreshToken models.RefreshToken) error {
	result := r.Db.Save(&refreshToken)
	return result.Error

}
func (r *refreshTokenRepositoryImpl) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	result := r.Db.Where("token = ?", token).First(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}
func (r *refreshTokenRepositoryImpl) Delete(token string) error {
	var refreshToken models.RefreshToken
	result := r.Db.Where("token = ?", token).Delete(&refreshToken)
	return result.Error
}
