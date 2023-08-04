package jwt_test

import (
	"testing"
	"time"

	"github.com/dmitrymomot/go-app/pkg/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTInteractor(t *testing.T) {
	// Test case 1: generate valid token and verify it
	t.Run("generate valid token and verify it", func(t *testing.T) {
		signingKey := []byte("secret")
		jwti := jwt.NewInteractor(signingKey, "test", time.Hour)
		tokenString, err := jwti.GenerateToken("token-id", "user-id", 0, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		claims, err := jwti.ValidateToken(tokenString, "test-audience")
		require.NoError(t, err)
		require.Equal(t, "token-id", claims.ID)
		require.Equal(t, "user-id", claims.Subject)
		require.Equal(t, "test-audience", claims.Audience[0])
		require.Equal(t, "test", claims.Issuer)
		require.NotEmpty(t, claims.ExpiresAt)
		require.Equal(t, claims.ExpiresAt.Time, claims.IssuedAt.Time.Add(time.Hour))
		require.NotEmpty(t, claims.IssuedAt)
		require.NotEmpty(t, claims.NotBefore)
	})

	// Test case 2: invalid token: expired
	t.Run("invalid token: expired", func(t *testing.T) {
		signingKey := []byte("secret")
		jwti := jwt.NewInteractor(signingKey, "test", 0)
		tokenString, err := jwti.GenerateToken("token-id", "user-id", 1, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		claims, err := jwti.ValidateToken(tokenString, "test-audience")
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenExpired)
		require.Nil(t, claims)
	})

	// Test case 3: invalid token: invalid audience
	t.Run("invalid token: invalid audience", func(t *testing.T) {
		signingKey := []byte("secret")
		jwti := jwt.NewInteractor(signingKey, "test", 0)
		tokenString, err := jwti.GenerateToken("token-id", "user-id", 1, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		claims, err := jwti.ValidateToken(tokenString, "wrong-audience")
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenInvalidAudience)
		require.Nil(t, claims)
	})

	// Test case 4: invalid token: invalid signature
	t.Run("invalid token: invalid signature", func(t *testing.T) {
		signingKey := []byte("secret")
		jwti := jwt.NewInteractor(signingKey, "test", 0)
		tokenString, err := jwti.GenerateToken("token-id", "user-id", 1, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		jwti = jwt.NewInteractor([]byte("secret2"), "test", 0)
		claims, err := jwti.ValidateToken(tokenString, "test-audience")
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
		require.Nil(t, claims)
	})

	// Test case 5: invalid token: invalid issuer
	t.Run("invalid token: invalid issuer", func(t *testing.T) {
		signingKey := []byte("secret")
		jwti := jwt.NewInteractor(signingKey, "test", 0)
		tokenString, err := jwti.GenerateToken("token-id", "user-id", 1, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		jwti = jwt.NewInteractor(signingKey, "test2", 0)
		claims, err := jwti.ValidateToken(tokenString, "test-audience")
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenInvalidIssuer)
		require.Nil(t, claims)
	})

	// Test case 6: generate valid token with anonymous user and verify it
	t.Run("generate valid token with anonymous user and verify it", func(t *testing.T) {
		signingKey := []byte("secret")
		jwti := jwt.NewInteractor(signingKey, "test", time.Hour)
		tokenString, err := jwti.GenerateToken("", "", 0, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		claims, err := jwti.ValidateToken(tokenString, "test-audience")
		require.NoError(t, err)
		require.NotEmpty(t, claims.ID)
		require.Equal(t, "anonymous", claims.Subject)
	})
}
