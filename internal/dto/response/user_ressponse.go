package response

type UserResponse struct {
	FullName string `validate:"required,min=2,max=100" json:"fullname"`
	Age      int    `validate:"required,min=2,max=100" json:"age"`
}

type LoginResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
