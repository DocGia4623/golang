package repository

import (
	"errors"
	"fmt"
	"testwire/helper"
	"testwire/internal/dto/request"
	"testwire/internal/models"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: db}
}

func (u *UserRepositoryImpl) Save(user request.CreateUserRequest) error {
	newUser := models.User{
		FullName: user.FullName,
		UserName: user.UserName,
		Password: user.Password, // Nên hash password trước khi lưu
		Age:      uint(user.Age),
	}
	result := u.Db.Create(&newUser)
	if result.Error != nil {
		return fmt.Errorf("failed to save user: %w", result.Error)
	}
	return nil
}

func (u *UserRepositoryImpl) FindById(userId int) (models.User, error) {
	var existingUser models.User
	err := u.Db.First(&existingUser, userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, err
	}
	return existingUser, nil
}

func (u *UserRepositoryImpl) Update(user request.UpdateUserRequest, userId int) error {
	existingUser, err := u.FindById(userId)
	if err != nil {
		return err
	}

	// Tạo map chứa các field cần cập nhật
	updateData := map[string]interface{}{}
	if user.FullName != "" {
		updateData["fullname"] = user.FullName
	}
	if user.Age > 0 {
		updateData["age"] = user.Age
	}
	if user.Password != "" {
		updateData["password"] = user.Password // chưa hash
	}

	// Cập nhật chỉ những field có giá trị hợp lệ
	if err := u.Db.Model(&existingUser).Updates(updateData).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *UserRepositoryImpl) Delete(userId int) error {
	var user models.User
	result := u.Db.Where("id = ?", userId).Delete(&user)
	if result.Error != nil {
		return fmt.Errorf("unable to delete user: %w", result.Error)
	}
	return nil
}
func (u *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	result := u.Db.Preload("Roles").Find(&users)
	if result.Error != nil {
		return users, fmt.Errorf("error finding users: %w", result.Error)
	}
	return users, nil
}

func (u *UserRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := u.Db.Where("user_name = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Trả về nil nếu không tìm thấy bản ghi
		}
		helper.ErrorPanic(result.Error)
		return nil, result.Error // Trả về lỗi nếu có lỗi khác từ GORM
	}
	return &user, nil
}

func (u *UserRepositoryImpl) FindIfUserHasRole(userID uint, roles []models.Role) error {
	var count int64

	// 🔹 Trích xuất tên role từ danh sách roles
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// 🔹 Truy vấn kiểm tra User có Role không
	result := u.Db.Model(&models.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("users.id = ? AND roles.name IN ?", userID, roleNames).
		Select("COUNT(*)").Scan(&count)

	// 🔹 Kiểm tra lỗi query
	if result.Error != nil {
		return result.Error
	}

	// 🔹 Kiểm tra nếu không tìm thấy Role nào
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil // ✅ User có ít nhất một Role phù hợp
}

func (u *UserRepositoryImpl) AddRole(userId uint, roles []models.Role) error {
	// 🔹 1. Tìm user trong database
	var user models.User
	if err := u.Db.Preload("Roles").First(&user, userId).Error; err != nil {
		return err // Trả về lỗi nếu không tìm thấy user
	}

	// 🔹 2. Gán danh sách roles cho user
	if err := u.Db.Model(&user).Association("Roles").Append(roles); err != nil {
		return err
	}

	return nil // Thành công
}

func (u *UserRepositoryImpl) FindRole(roles []string) ([]models.Role, error) {
	var foundRoles []models.Role

	if err := u.Db.Where("name IN ?", roles).Find(&foundRoles).Error; err != nil {
		return nil, err
	}

	return foundRoles, nil
}
