package request

type SignUpRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
}

type SignInRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required,min=3"`
	Password        string `json:"password" binding:"required,min=8,max=100"`
}
