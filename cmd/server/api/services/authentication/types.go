package authentication

type userRegisterPayload struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type userLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type userForgotPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type userPasswordResetPayload struct {
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type googleLoginPayload struct {
	Email      string `json:"email" validate:"required,email"`
	GivenName  string `json:"given_name" validate:"required"`
	FamilyName string `json:"family_name" validate:"required"`
}
