package models

import "regexp"

// Constants for error/success messages
const (
	CHOOSE_BETTER_PASSCODE_MESSAGE = "choose a better passcode"
	INVALID_INPUT_MESSAGE          = "Invalid Input."
	INVALID_CREDENTIALS            = "Invalid credentials"
	PERSONAL_USER_CREATED          = "Successfully registered."
	INVALID_TOKEN_MESSAGE          = "Invalid token"
	LOGIN_SUCCESSFUL               = "User successfully logged In"
	USER_NOT_FOUND                 = "user info not found"
	USER_FETCHED_SUCCESSFULLY      = "User fetched successfully"
	TOKEN_EXPIRED_MESSAGE          = "Token Expired"
	INTERNAL_SERVER_ERROR_MESSAGE  = "Internal Server error"
	USER_ALREADY_EXIST             = "User already exists with this email"
)

// Regular expressions for validation
var (
	// nameRegex ensures the name consists of alphanumeric characters and is between 2 and 26 characters long
	nameRegex = regexp.MustCompile("^[a-zA-Z0-9]{2,26}$")

	// emailRegex validates the email format
	emailRegex = regexp.MustCompile("^([a-zA-Z0-9!&'*+\\/=?^_`{|}-]+(?:\\.[a-zA-Z0-9!#$%&'*+\\/=?^_`{|}~-]+)*@(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?\\.)+[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?)$")
)
