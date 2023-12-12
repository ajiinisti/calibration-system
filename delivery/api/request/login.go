package request

type Login struct {
	Email    string
	Password string
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// ? ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
