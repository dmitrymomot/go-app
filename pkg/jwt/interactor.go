package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Interactor is the interface that provides the methods that are used by the
// application layer.
type Interactor interface {
	// GenerateToken generates a new JWT token for the given token ID and subject.
	// The token will expire after the given TTL. The token will be valid for the
	// given audience(s). Subject can be empty, in which case the token will be
	// generated for an anonymous user.
	GenerateToken(id, subject string, ttl time.Duration, audience ...string) (string, error)

	// ValidateToken validates the given JWT token and returns the claims if the
	// token is valid. The token will be validated for the given audience.
	ValidateToken(token string, audience string) (*Claims, error)
}

// interactor is the implementation of the Interactor interface.
type interactor struct {
	signingKey []byte
	issuer     string
	ttl        time.Duration
}

// NewInteractor returns a new Interactor instance.
func NewInteractor(signingKey []byte, issuer string, ttl time.Duration) Interactor {
	if ttl == 0 {
		ttl = time.Hour
	}
	if issuer == "" {
		issuer = "go-app/pkg/jwt"
	}
	return &interactor{
		signingKey: signingKey,
		issuer:     issuer,
		ttl:        ttl,
	}
}

// GenerateToken generates a new JWT token for the given token ID and subject.
// The token will expire after the given TTL. The token will be valid for the
// given audience(s). Subject can be empty, in which case the token will be
// generated for an anonymous user.
func (i *interactor) GenerateToken(id, subject string, ttl time.Duration, audience ...string) (string, error) {
	if ttl == 0 {
		ttl = i.ttl
	}
	if id == "" {
		id = uuid.New().String()
	}
	if subject == "" {
		subject = "anonymous"
	}

	// Create the claims for the token with standard claims.
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    i.issuer,
			Subject:   subject,
			ID:        id,
			Audience:  audience,
		},
	}

	// Create the signed token string with the claims.
	tokenString, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(i.signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the given JWT token and returns the claims if the
// token is valid. The token will be validated for the given audience.
func (i *interactor) ValidateToken(tokenString string, audience string) (*Claims, error) {
	// Parse the token string into a token object.
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return i.signingKey, nil
		},
		jwt.WithAudience(audience),
		jwt.WithIssuer(i.issuer),
	)
	if err != nil {
		return nil, err
	}

	// Check if the token and claims are valid.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
