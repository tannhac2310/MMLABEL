package jwtutil

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const (
	secretKey = "secret_key"
	jwtAud    = "some_aud"
	jwtIss    = "some_iss"
	exp       = time.Second * 3600
)

func TestTokenGenerator(t *testing.T) {
	var (
		claims = &Claims{
			Roles: []string{"root"},
		}
		err                    error
		generator              TokenGenerator
		userId                 = ksuid.New().String()
		jwtToken, refreshToken string
	)

	t.Run("New Token Generator", func(t *testing.T) {
		generator, err = NewTokenGenerator(secretKey, jwtAud, jwtIss, exp)
		assert.Nil(t, err)
		assert.NotNil(t, generator)
	})

	t.Run("Generate Token", func(t *testing.T) {
		jwtToken, refreshToken, err = generator.Encode(userId, claims)
		assert.Nil(t, err)
		assert.NotEmpty(t, jwtToken)
		assert.NotEmpty(t, refreshToken)
	})

	t.Run("Decode JWT expired", func(t *testing.T) {
		generatorTest, err := NewTokenGenerator(secretKey, jwtAud, jwtIss, 1)
		assert.Nil(t, err)
		_jwtToken, _, err := generatorTest.Encode(userId, claims)
		assert.Nil(t, err)

		time.Sleep(2 * time.Second)
		subject, vars, err := generatorTest.DecodeJwt(_jwtToken)
		assert.Equal(t, "exp not satisfied", err.Error())
		assert.Equal(t, "", subject)
		assert.Nil(t, vars)
	})

	t.Run("Decode JWT", func(t *testing.T) {
		subject, vars, err := generator.DecodeJwt(jwtToken)
		assert.Nil(t, err)
		assert.Equal(t, userId, subject)
		assert.Equal(t, "root", vars.Roles[0])
	})

	t.Run("Decode Refresh Token", func(t *testing.T) {
		subject, err := generator.DecodeRefreshToken(refreshToken)
		assert.Nil(t, err)
		assert.Equal(t, userId, subject)
	})

	t.Run("Decode JWT with Refresh Token", func(t *testing.T) {
		subject, vars, err := generator.DecodeJwt(refreshToken)
		assert.NotNil(t, err)
		assert.Nil(t, vars)
		assert.Equal(t, "failed to verify jws signature: failed to verify message: failed to match hmac signature", err.Error())
		assert.Empty(t, subject)
	})

	t.Run("Decode Refresh Token with JWT", func(t *testing.T) {
		subject, err := generator.DecodeRefreshToken(jwtToken)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to verify jws signature: failed to verify message: failed to match hmac signature", err.Error())
		assert.Empty(t, subject)
	})
}
