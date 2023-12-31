package repository

import "test/models"

type Store interface {
	AddUser(user models.UserSignUpRequest) error
	GetUserByEmail(email string) (models.User, error)
	AuthenticateUser(loginDetail models.UserSignInRequest) (bool, error)
}
