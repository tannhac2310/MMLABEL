package auth

import (
	"context"
	"fmt"
)

func (a *authService) RefreshToken(ctx context.Context, token string) (*LoginResult, error) {
	userID, err := a.jwtGenerator.DecodeRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("a.jwtGenerator.DecodeRefreshToken: %w", err)
	}

	return a.buildLoginResult(ctx, userID)
}
