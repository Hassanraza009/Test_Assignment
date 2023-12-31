package models

// UserSignUpRequest represents the structure of a user signup request.
type UserSignUpRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	HashPassword string // HashPassword is used to store the hashed password securely.
}

// User represents the structure of a user entity.
type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// UserSignInRequest represents the structure of a user sign-in request.
type UserSignInRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	HashPassword string // HashPassword is used to store the hashed password securely.
}
