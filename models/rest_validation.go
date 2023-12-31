package models

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// validates the user signup request
func (s *UserSignUpRequest) Validate() error {

	invalidInputError := fmt.Errorf("invalid data")

	if areAllSame(s.Password) {
		return fmt.Errorf(CHOOSE_BETTER_PASSCODE_MESSAGE)
	}

	if areDigitsInSequense(s.Password) {
		return fmt.Errorf(CHOOSE_BETTER_PASSCODE_MESSAGE)
	}

	if !isNameValid(s.FirstName) {
		return invalidInputError
	}
	if !isNameValid(s.LastName) {
		return invalidInputError
	}
	if !isEmailValid(s.Email) {
		return invalidInputError
	}
	return nil
}

func (s *UserSignInRequest) Validate() error {

	invalidInputError := fmt.Errorf("invalid data")
	if !isEmailValid(s.Email) {
		return invalidInputError
	}
	return nil
}

// check if all the characters in the string are same
func areAllSame(val string) bool {
	var allSame bool
	lengthOfText := len(val)
	firstVal := []rune(val)[0]
	reg := regexp.MustCompile(string(firstVal))
	content := []byte(val)

	occurences := reg.FindAllIndex(content, lengthOfText)
	if occurences != nil && len(occurences) == lengthOfText {
		allSame = true
	}

	return allSame
}

// check if we have consective numbers
func areDigitsInSequense(value string) bool {
	var found bool
	var allDigits = "0123456789"

	found = strings.Contains(allDigits, value)

	return found
}

// isNameValid checks if the passcode provided passes the required structure
func isNameValid(e string) bool {

	// ISO IEC 7813
	if len(e) > 26 || len(e) < 2 {
		return false
	}
	if !nameRegex.MatchString(e) {
		return false
	}
	return true
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
