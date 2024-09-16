package request

type SignUp struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string `json:"last_name" validate:"required,min=3,max=32"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	// Future use for Oauth2.
	// Provider string `json:"provider"`
}
