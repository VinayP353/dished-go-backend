package models

type CreateCookRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateCookRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=3"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Status   *string `json:"status,omitempty"`
}

type CreateCookProfileRequest struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	PreferredName string `json:"preferred_name"`
	Address       string `json:"address"`
	Description   string `json:"description"`
}

type UpdateCookProfileRequest struct {
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	PreferredName *string `json:"preferred_name,omitempty"`
	Address       *string `json:"address,omitempty"`
	Description   *string `json:"description,omitempty"`
	Verified      *bool   `json:"verified,omitempty"`
}
