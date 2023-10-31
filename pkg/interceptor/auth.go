package interceptor

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/jwtutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const AllJwtTokenValidAfter = "all:jwt_token_valid_after" //nolint:gosec

func Auth(jwt jwtutil.TokenGenerator, e casbin.IEnforcer, redisDB redis.Cmdable, decider func(c *gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !decider(c) {
			c.Next()
			return
		}

		token := c.Query("token")
		if token == "" {
			authorization := c.GetHeader("Authorization")
			parts := strings.Split(authorization, " ")
			if len(parts) != 2 {
				_ = c.Error(apperror.ErrUnauthenticated.WithDebugMessage("wrong authorization token format"))
				c.Abort()
				return
			}

			if parts[0] != "Bearer" {
				_ = c.Error(apperror.ErrUnauthenticated.WithDebugMessage("wrong authorization token format, must start with Bearer"))
				c.Abort()
				return
			}

			token = parts[1]
		}

		jwtToken, userID, claims, err := jwt.DecodeJwt(token)
		if err != nil {
			_ = c.Error(apperror.ErrUnauthenticated.WithDebugMessage(err.Error()))
			c.Abort()
			return
		}
		var expTime int64
		err = redisDB.Get(AllJwtTokenValidAfter).Scan(&expTime)
		if err != nil {
			ctxzap.Extract(c).Warn("err get cache", zap.String("key", AllJwtTokenValidAfter), zap.String("userId", userID), zap.Error(err))
		}
		if expTime > jwtToken.IssuedAt().Unix() {
			_ = c.Error(apperror.ErrUnauthenticated.WithDebugMessage("token is revoked"))
			c.Abort()
			return
		}

		path := c.Request.URL.Path

		isAdmin := false
		canAccess := false
		rolesAndGroups := append(claims.Roles, claims.Groups...)
		for _, p := range rolesAndGroups {
			if p == model.UserRoleRoot {
				isAdmin = true
				canAccess = true
				break
			}

			allowed, err := e.Enforce(p, path)
			if err != nil {
				_ = c.Error(fmt.Errorf("e.Enforce: %w", err))
				c.Abort()
				return
			}

			if allowed {
				canAccess = true
				break
			}
		}

		if !canAccess {
			_ = c.Error(apperror.ErrPermissionDenied.WithDebugMessage("can not access this endpoint " + path))
			c.Abort()
			return
		}
		// Todo verify trong table permission

		c.Set("isAdmin", isAdmin)
		c.Set("userID", userID)
		c.Next()
	}
}
func UserIDFromCtx(c *gin.Context) string {
	return c.GetString("userID")
}

func IsAdmin(c *gin.Context) bool {
	return c.GetBool("isAdmin")
}

func OwnerIDIfNotAdmin(c *gin.Context) string {
	if IsAdmin(c) {
		return ""
	}

	return UserIDFromCtx(c)
}

func ForceTokenInValid(ctx context.Context, cfg *configs.BaseConfig, redisDB redis.Cmdable) {
	//   todo uncomment later
	// err := multierr.Combine(
	//	redisDB.Set(AllJwtTokenValidAfter, time.Now().Unix(), cfg.JWT.Expiry).Err(),
	//)
	//if err != nil {
	//	ctxzap.Extract(ctx).Warn("err store cache", zap.String("key", AllJwtTokenValidAfter), zap.Error(err))
	//}
}
