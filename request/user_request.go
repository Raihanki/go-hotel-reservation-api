package request

type RegisterRequest struct {
	Name                 string `json:"name" form:"name" validate:"required,min=3,max=100"`
	Email                string `json:"email" form:"email" validate:"required,email"`
	Password             string `json:"password" form:"password" validate:"required,min=8,max=100"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" validate:"required,min=8,max=100,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}
