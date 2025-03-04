//go:build wireinject
// +build wireinject

package wire

import (
	"testwire/config"
	"testwire/internal/controller"
	"testwire/internal/middleware"
	"testwire/internal/repository"
	"testwire/internal/services"

	"github.com/google/wire"
)

// Inject các dependency
var AppSet = wire.NewSet(
	config.LoadConfig,
	config.ConnectDB,           // Kết nối DB
	RepositorySet,              // Inject Repository
	MiddlerwareSet,             // Inject Middleware
	ServiceSet,                 // Inject Service
	ControllerSet,              // Inject Controller
	wire.Struct(new(App), "*"), // Tự động inject vào struct Ap
)

type App struct {
	AuthController    *controller.AuthenticationController
	AuthService       services.AuthenticationService
	AuthRepo          repository.UserRepository
	UserController    *controller.UserController
	Middleware        *middleware.Middleware
	ProductController *controller.ProductController
	OrderController   *controller.OrderController
}

// InitializeUserService khởi tạo UserService tự động
func InitializeApp() (*App, error) {
	wire.Build(AppSet)
	return nil, nil
}
