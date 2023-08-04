package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Predefined errors.
var (
	ErrInvalidToken = errors.New("invalid or expired token")

	ErrTokenExpired              = jwt.ErrTokenExpired
	ErrTokenInvalidAudience      = jwt.ErrTokenInvalidAudience
	ErrInvalidKey                = jwt.ErrInvalidKey
	ErrInvalidKeyType            = jwt.ErrInvalidKeyType
	ErrHashUnavailable           = jwt.ErrHashUnavailable
	ErrTokenMalformed            = jwt.ErrTokenMalformed
	ErrTokenUnverifiable         = jwt.ErrTokenUnverifiable
	ErrTokenSignatureInvalid     = jwt.ErrTokenSignatureInvalid
	ErrTokenRequiredClaimMissing = jwt.ErrTokenRequiredClaimMissing
	ErrTokenUsedBeforeIssued     = jwt.ErrTokenUsedBeforeIssued
	ErrTokenInvalidIssuer        = jwt.ErrTokenInvalidIssuer
	ErrTokenInvalidSubject       = jwt.ErrTokenInvalidSubject
	ErrTokenNotValidYet          = jwt.ErrTokenNotValidYet
	ErrTokenInvalidId            = jwt.ErrTokenInvalidId
	ErrTokenInvalidClaims        = jwt.ErrTokenInvalidClaims
	ErrInvalidType               = jwt.ErrInvalidType
)
