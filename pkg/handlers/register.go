package handlers

// RegisterRequest represents the request payload for user registration.
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// RegisterResponse represents the response payload for a successful registration.
type RegisterResponse struct {
	Message string `json:"message"`
}
