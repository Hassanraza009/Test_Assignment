package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"test/logger"
	"test/models"
	"test/repository"
)

type UserService interface {
	CreateUser(signupRequest models.UserSignUpRequest) error
	HashPassword(password string) (string, error)
	Login(loginDetails models.UserSignInRequest) (string, error)
	GetUserByEmail(email string) (*models.User, error)
}
type userService struct {
	logger     logger.Logger
	jwtService JWTService
	store      repository.Store
}

func NewUserService(store repository.Store, logger logger.Logger, jwtService JWTService) UserService {
	return &userService{
		jwtService: jwtService,
		store:      store,
		logger:     logger,
	}
}
func (u *userService) CreateUser(signupRequest models.UserSignUpRequest) error {

	existingUser, err := u.store.GetUserByEmail(signupRequest.Email)
	if err != nil {
		return err
	}

	// If the user with the given email already exists, return an error
	if existingUser.Id != 0 {
		return fmt.Errorf("user with email %s already exists", signupRequest.Email)
	}

	signupRequest.HashPassword, err = u.HashPassword(signupRequest.Password)
	if err != nil {
		return err
	}
	// If the user does not exist, proceed with creating the user
	err = u.store.AddUser(signupRequest)
	if err != nil {
		return err
	}
	return nil
}

// Login authenticates a user based on the provided login details,
// generates a JWT token for the authenticated user, and returns the token.
func (u *userService) Login(loginDetails models.UserSignInRequest) (string, error) {
	// Hash the user's password for secure storage and comparison
	var err error
	loginDetails.HashPassword, err = u.HashPassword(loginDetails.Password)
	if err != nil {
		return "", err
	}

	// Authenticate the user based on the provided login details
	status, err := u.store.AuthenticateUser(loginDetails)
	if err != nil {
		return "", err
	}
	if !status {
		return "", fmt.Errorf("invalid credentials")
	}

	// Retrieve user details after successful authentication
	userDetail, err := u.store.GetUserByEmail(loginDetails.Email)
	if err != nil {
		return "", err
	}

	// Generate a JWT token for the authenticated user
	token, err := u.jwtService.CreateLoginToken(userDetail)
	if err != nil {
		return "", err
	}
	// Return the generated token
	return token, nil
}

func (u *userService) HashPassword(password string) (string, error) {
	hasher := sha256.New()

	// Write the password bytes to the hasher
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}

	// Get the hashed password as a byte slice
	hashedPassword := hasher.Sum(nil)

	// Convert the hashed password to a hex-encoded string
	hashedPasswordString := hex.EncodeToString(hashedPassword)

	return hashedPasswordString, nil
}
func (u *userService) GetUserByEmail(email string) (*models.User, error) {
	userDetail, err := u.store.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &userDetail, nil
}
