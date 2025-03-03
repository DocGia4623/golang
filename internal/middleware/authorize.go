package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testwire/config"
	"testwire/helper"
	"testwire/internal/repository"
	"testwire/utils"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	UserRepo       repository.UserRepository
	PermissionRepo repository.PermissionRepository
}

func NewMiddleware(userRepo repository.UserRepository, permissionRepo repository.PermissionRepository) *Middleware {
	return &Middleware{
		UserRepo:       userRepo,
		PermissionRepo: permissionRepo,
	}
}

func (m *Middleware) AuthorizeRole(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get token
		var token string
		authHeader := ctx.GetHeader("Authorization")
		fields := strings.Fields(authHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		}
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"status": "fail", "message": "Missing token"})
			return
		}
		config, _ := config.LoadConfig()
		sub, err := utils.ValidateAccessToken(token, config.AccessTokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		id, err_id := strconv.Atoi(fmt.Sprint(sub))
		helper.ErrorPanic(err_id)
		_, err = m.UserRepo.FindById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "User not found"})
			return
		}

		permissionModel, err := m.PermissionRepo.FindIfExist(permission)
		if err != nil || permissionModel == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission not found"})
			return
		}

		roles, err := m.PermissionRepo.FindRoleBelong(permissionModel.Name)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "error find role"})
			return
		}
		if len(roles) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "empty role"})
			return
		}
		idUser := uint(id)
		// Check if user has role
		err = m.UserRepo.FindIfUserHasRole(idUser, roles)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
			return
		}
	}
}
