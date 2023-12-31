package models

// DecodeJWTClaims represents the decoded claims from a JWT token.
type DecodeJWTClaims struct {
	ExpiredAt float64 `json:"expiredAt"`
	UID       float64 `json:"uid"`
	Email     string  `json:"email"`
}
