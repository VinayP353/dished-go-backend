package models

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=8"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Chef  *Chef  `json:"chef"`
	Token string `json:"token,omitempty"`
}

type UsernameResponse struct {
	Username string `json:"username"`
}
