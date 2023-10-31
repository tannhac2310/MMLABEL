package transportutil

import (
	"github.com/casbin/casbin/v2"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/jwtutil"
)

var skipLogs map[string]struct{} = map[string]struct{}{
	"/metrics": {},
	"/healthz": {},
}

var publicEndpoint = map[string]struct{}{
	"/hydra/api-docs":  {},
	"/gezu/api-docs":   {},
	"/aurora/api-docs": {},
	"/favicon.ico":     {},
	"/metrics":         {},
	"/healthz":         {},
}

func RegisterPublicEndpoint(path string) {
	publicEndpoint[path] = struct{}{}
}

func InitGinEngine(
	cfg *configs.BaseConfig,
	jwt jwtutil.TokenGenerator,
	e casbin.IEnforcer,
	zapLogger *zap.Logger,
	ext cockroach.Ext,
	redisDB redis.Cmdable,
) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if !cfg.Debug {
		r.Use(ginzap.RecoveryWithZap(zapLogger, true))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// inject db connection to ctx
	r.Use(func(c *gin.Context) {
		ctx := cockroach.ContextWithDB(c.Request.Context(), ext)
		c.Request = c.Request.WithContext(ctx)
	})
	r.Use(otelgin.Middleware(cfg.Name))
	r.Use(HandleError)
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "DeviceID", "Accept-Language"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))
	r.Use(requestid.New())

	r.Use(interceptor.Ginzap(zapLogger, func(c *gin.Context) bool {
		_, ok := skipLogs[c.Request.URL.Path]
		return !(ok || strings.Contains(c.Request.URL.Path, "-docs"))
	}, cfg.Logger.LogReq, cfg.Logger.LogResp))

	r.Use(interceptor.Auth(jwt, e, redisDB, func(c *gin.Context) bool {
		_, ok := publicEndpoint[c.Request.URL.Path]
		return !(ok || strings.Contains(c.Request.URL.Path, "-docs"))
	}))

	return r
}
