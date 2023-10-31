package jwtutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type Claims struct {
	Groups []string `json:"groups"`
	Roles  []string `json:"roles"`
}

type TokenGenerator interface {
	Encode(userID string, c *Claims) (jwtToken, refreshToken string, err error)
	DecodeJwt(token string) (t jwt.Token, userID string, claims *Claims, err error)
	DecodeRefreshToken(refreshToken string) (userID string, err error)
}

type tokenGenerator struct {
	jwtEncryptionKey         []byte
	refreshEncryptionKey     []byte
	jwtAudience, jwtIssuer   string
	jwtExpiry, refreshExpiry time.Duration
}

func NewTokenGenerator(jwtEncryptionKey string, jwtAudience, jwtIssuer string, jwtExpiry time.Duration) (TokenGenerator, error) {
	if len(jwtEncryptionKey) == 0 {
		return nil, fmt.Errorf("jwtEncryptionKey should not empty")
	}

	if jwtAudience == "" {
		return nil, fmt.Errorf("jwtAudience should not empty")
	}

	if jwtIssuer == "" {
		return nil, fmt.Errorf("jwtIssuer should not empty")
	}

	if jwtExpiry == 0 {
		return nil, fmt.Errorf("jwtExpiry should not zero")
	}

	return &tokenGenerator{
		jwtEncryptionKey:     []byte(jwtEncryptionKey),
		refreshEncryptionKey: []byte("refresh_token" + jwtEncryptionKey),
		jwtAudience:          jwtAudience,
		jwtIssuer:            jwtIssuer,
		jwtExpiry:            jwtExpiry,
		refreshExpiry:        30 * 7 * 24 * time.Hour,
	}, nil
}

func (rcv *tokenGenerator) Encode(userID string, claims *Claims) (string, string, error) {
	var (
		now          = time.Now().UTC()
		jwtToken     = jwt.New()
		refreshToken = jwt.New()
	)

	_ = jwtToken.Set(jwt.SubjectKey, userID)
	_ = jwtToken.Set(jwt.AudienceKey, rcv.jwtAudience)
	_ = jwtToken.Set(jwt.ExpirationKey, now.Add(rcv.jwtExpiry))
	_ = jwtToken.Set(jwt.IssuedAtKey, now)
	_ = jwtToken.Set(jwt.IssuerKey, rcv.jwtIssuer+"/"+rcv.jwtAudience)
	if claims != nil {
		_ = jwtToken.Set("claims", claims)
	}

	jwtSigned, err := jwt.Sign(jwtToken, jwa.HS256, rcv.jwtEncryptionKey)
	if err != nil {
		return "", "", errors.Wrap(err, "jwtToken.Sign")
	}

	_ = refreshToken.Set(jwt.SubjectKey, userID)
	_ = refreshToken.Set(jwt.AudienceKey, rcv.jwtAudience)
	_ = refreshToken.Set(jwt.ExpirationKey, now.Add(rcv.refreshExpiry))
	_ = refreshToken.Set(jwt.IssuedAtKey, now)
	_ = refreshToken.Set(jwt.IssuerKey, rcv.jwtIssuer+"/"+rcv.jwtAudience)

	refreshTokenSigned, err := jwt.Sign(refreshToken, jwa.HS256, rcv.refreshEncryptionKey)
	if err != nil {
		return "", "", errors.Wrap(err, "refreshToken.Sign")
	}

	return string(jwtSigned), string(refreshTokenSigned), nil
}

func (rcv *tokenGenerator) DecodeJwt(jwtToken string) (t jwt.Token, userID string, claims *Claims, err error) {
	token, err := jwt.ParseString(
		jwtToken,
		jwt.WithVerify(jwa.HS256, rcv.jwtEncryptionKey),
		jwt.WithAudience(rcv.jwtAudience),
		jwt.WithIssuer(rcv.jwtIssuer+"/"+rcv.jwtAudience),
		jwt.WithClock(jwt.ClockFunc(func() time.Time {
			return time.Now().UTC()
		})),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, "", nil, err
	}

	vars, ok := token.Get("claims")
	if !ok {
		return nil, "", nil, fmt.Errorf("missing `claims` in token")
	}

	data, _ := json.Marshal(vars)
	claims = &Claims{}
	_ = json.Unmarshal(data, claims)

	return token, token.Subject(), claims, nil
}

func (rcv *tokenGenerator) DecodeRefreshToken(refreshToken string) (userID string, err error) {
	token, err := jwt.ParseString(
		refreshToken,
		jwt.WithVerify(jwa.HS256, rcv.refreshEncryptionKey),
		jwt.WithAudience(rcv.jwtAudience),
		jwt.WithIssuer(rcv.jwtIssuer+"/"+rcv.jwtAudience),
		jwt.WithClock(jwt.ClockFunc(func() time.Time {
			return time.Now().UTC()
		})),
		jwt.WithValidate(true),
	)
	if err != nil {
		return "", err
	}

	return token.Subject(), nil
}
