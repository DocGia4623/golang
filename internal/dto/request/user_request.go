package request

type CreateUserRequest struct {
	UserName string `validate:"required,min=2,max=100" json:"username"`
	FullName string `validate:"required,min=2,max=100" json:"fullname"`
	Age      int    `validate:"required,min=2,max=100" json:"age"`
	Password string `validate:"required,min=2,max=100" json:"password"`
}
type UpdateUserRequest struct {
	FullName string `validate:"required,min=2,max=100" json:"fullname"`
	Age      int    `validate:"required,min=2,max=100" json:"age"`
	Password string `validate:"required,min=2,max=100" json:"password"`
}

type LoginRequest struct {
	UserName string `validate:"required,min=2,max=100" json:"username"`
	Password string `validate:"required,min=2,max=100" json:"password"`
}

type AddRoleRequest struct {
	UserID uint     `json:"user_id"`
	Roles  []string `json:"roles"`
}
