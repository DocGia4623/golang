package services

import (
	"context"
	"errors"
	"fmt"
	"testwire/config"
	"testwire/helper"
	"testwire/internal/dto/request"
	"testwire/internal/repository"
	"testwire/utils"
	"time"

	"gorm.io/gorm"
)

type AuthenticationService interface {
	Register(request.CreateUserRequest) error
	Login(string, string) (string, string, error)
	Logout(context.Context, string, string) error
}

type AuthenticationServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewAuthenticationServiceImpl(userRepo repository.UserRepository) AuthenticationService {
	return &AuthenticationServiceImpl{UserRepo: userRepo}
}

func (a *AuthenticationServiceImpl) Register(userRequest request.CreateUserRequest) error {
	// Kiểm tra xem username đã tồn tại chưa
	user, userErr := a.UserRepo.FindByUsername(userRequest.UserName)
	if userErr != nil && !errors.Is(userErr, gorm.ErrRecordNotFound) {
		return userErr
	}
	if user != nil {
		// Nếu user đã tồn tại trong cơ sở dữ liệu
		return errors.New("username already exists")
	}
	// Hash mật khẩu
	hashedPassword, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		return err // Nếu có lỗi khi hash mật khẩu, trả về lỗi
	}
	userRequest.Password = hashedPassword

	// Lưu user vào database
	err = a.UserRepo.Save(userRequest)
	if err != nil {
		return fmt.Errorf("error while saving to db: %w", err)
	}
	return nil // Trả về nil nếu không có lỗi
}
func (a *AuthenticationServiceImpl) Login(username string, password string) (string, string, error) {
	// Find username in the database
	login_user, user_err := a.UserRepo.FindByUsername(username)
	if user_err != nil {
		return "", "", errors.New("invalid username or password")
	}
	if login_user == nil {
		return "", "", errors.New("invalid username or password") // Xử lý khi user không tồn tại
	}
	//Verify password
	verify_err := utils.VerifyPassword(login_user.Password, password)
	if verify_err != nil {
		return "", "", errors.New("invalid username or password")
	}

	conf, err := config.LoadConfig()
	if err != nil {
		return "", "", errors.New("cant load config")
	}

	// Generate access token
	accessToken, err_token := utils.GenerateAccessToken(conf.AccessTokenExpiresIn, login_user.ID, conf.AccessTokenSecret)
	helper.ErrorPanic(err_token)

	// Generate refresh token
	refreshToken, err_refresh := utils.GenerateAccessToken(conf.RefreshTokenExpiresIn, login_user.ID, conf.RefreshTokenSecret)
	helper.ErrorPanic(err_refresh)

	return refreshToken, accessToken, nil
}

func (a *AuthenticationServiceImpl) Logout(ctx context.Context, refreshToken string, accessToken string) error {

	// Lưu token vào redis
	expiration := time.Hour
	err := config.RedisClient.Set(ctx, accessToken, "logout", expiration).Err()
	if err != nil {
		return err
	}

	// Xóa refresh token khỏi database
	RefreshTokenService := NewRefreshTokenServiceImpl(repository.NewRefreshTokenRepositoryImpl(config.DB))
	RefreshTokenService.DeleteToken(refreshToken)
	return nil
}
