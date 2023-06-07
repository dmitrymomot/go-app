package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims is a custom JWT claims.
type Claims struct {
	jwt.RegisteredClaims
}

// UserUUID returns a user ID as UUID.
func (c Claims) UserUUID() uuid.UUID {
	if c.Subject == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(c.Subject)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// TokenUUID returns a token ID as UUID.
func (c Claims) TokenUUID() uuid.UUID {
	if c.ID == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(c.ID)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// Audience returns an audience.
func (c Claims) Audiences() []string {
	if c.Audience == nil {
		return []string{}
	}
	return c.Audience
}

// AudienceExists checks if the audience exists.
func (c Claims) AudienceExists(audience string) bool {
	for _, a := range c.Audiences() {
		if a == audience {
			return true
		}
	}
	return false
}
