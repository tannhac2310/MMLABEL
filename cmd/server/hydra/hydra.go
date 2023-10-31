package hydra

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/nats-io/stan.go"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/controller"
	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/subscriptions"
	pkgConfig "mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/jwtutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/logger"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/nats"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/rbacutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/auth"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/casbinrule"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/upload"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/tracingutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func Run(ctx context.Context, configPath string) {
	cfg, err := configs.NewConfig(configPath)
	if err != nil {
		log.Panic("err load config:", err)
	}

	opts := []fx.Option{
		fx.Provide(func(cfg *pkgConfig.BaseConfig) *zap.Logger {
			return logger.NewZapLogger(cfg.Logger.LogLevel, cfg.Env == "local")
		}),
		// config
		fx.Provide(
			func() *configs.Config {
				return cfg
			},
			func(cfg *configs.Config) *pkgConfig.BaseConfig {
				return &cfg.BaseConfig
			},
			func(cfg *configs.Config) *configs.SmsBrand {
				return cfg.SmsBrand
			},
		),
		// http server
		fx.Provide(transportutil.InitGinEngine),
		// route group
		fx.Provide(func(r *gin.Engine) *gin.RouterGroup {
			g := r.Group("hydra")
			g.GET("api-docs", routeutil.ServingDocs)

			return g
		}),
		// redis
		fx.Provide(func(cfg *pkgConfig.BaseConfig) (redis.Cmdable, error) {
			c := redis.NewClient(&redis.Options{
				Addr:     cfg.Redis.Host,
				Password: string(cfg.Redis.Password),
				DB:       0, // use default DB
			})

			// Test
			_, err := c.Ping().Result()
			if err != nil {
				return nil, err
			}

			return c, nil
		}),
		// cockroach db
		fx.Provide(func(cfg *pkgConfig.BaseConfig, logger *zap.Logger) cockroach.Ext {
			return cockroach.NewConnectionPool(ctx, logger, string(cfg.CockroachDB.URI), cfg.Debug)
		}),
		// jwt
		fx.Provide(func(cfg *pkgConfig.BaseConfig) (jwtutil.TokenGenerator, error) {
			return jwtutil.NewTokenGenerator(
				cfg.JWT.EncryptionKey,
				cfg.JWT.Audience,
				cfg.JWT.Issuer,
				cfg.JWT.Expiry,
			)
		}),
		// casbin rbac
		fx.Provide(func(cfg *pkgConfig.BaseConfig) (casbin.IEnforcer, error) {
			return rbacutil.New(string(cfg.CockroachDB.URI))
		}),
		// firebase
		fx.Provide(func(cfg *configs.Config) (*firebaseAuth.Client, error) {
			opt := option.WithCredentialsFile(cfg.Firebase.ConfigPath)
			app, err := firebase.NewApp(ctx, nil, opt)
			if err != nil {
				return nil, fmt.Errorf("error initializing app: %v", err)
			}

			// Get an auth client from the firebase.App
			client, err := app.Auth(ctx)
			if err != nil {
				return nil, fmt.Errorf("error getting Auth client: %v", err)
			}

			return client, nil
		}),
		// monitoring
		fx.Provide(tracingutil.InitTelemetry),
		fx.Provide(func(cfg *configs.Config) (*profiler.Profiler, error) {
			if !cfg.RemoteProfiler.Enabled {
				return nil, nil
			}

			return profiler.Start(profiler.Config{
				ApplicationName: fmt.Sprintf("flamingo-group.%s", cfg.Name),
				ServerAddress:   cfg.RemoteProfiler.ProfilerURL,
			})
		}),
		// s3 storage
		fx.Provide(func(cfg *configs.Config) (*s3manager.Uploader, error) {
			sess, err := session.NewSession(&aws.Config{
				Region:           aws.String(cfg.S3Storage.Region),
				Credentials:      credentials.NewStaticCredentials(cfg.S3Storage.AccessKey, cfg.S3Storage.SecretKey, ""),
				Endpoint:         aws.String(cfg.S3Storage.Endpoint),
				S3ForcePathStyle: aws.Bool(true),
			})
			if err != nil {
				return nil, err
			}

			return s3manager.NewUploader(sess), nil
		}),
		// repositories
		fx.Provide(
			repository.NewUserRepo,
			repository.NewUserNamePasswordRepo,
			repository.NewUserRoleRepo,
			repository.NewUserGroupRepo,
			repository.NewGroupRepo,
			repository.NewRoleRepo,
			repository.NewUserFirebaseRepo,
			repository.NewUserFCMTokenRepo,
			repository.NewUserNotificationRepo,
			repository.NewCasbinRuleRepo,
		),
		// services
		fx.Provide(
			auth.NewService,
			user.NewService,
			group.NewService,
			role.NewService,
			upload.NewService,
			casbinrule.NewService,
		),
		// nats streaming
		fx.Provide(func(cfg *pkgConfig.BaseConfig, zapLogger *zap.Logger) (nats.BusFactory, error) {
			hostName, _ := os.Hostname()
			bus, err := nats.NewBusFactory(zapLogger, cfg.Nats.ClusterID, cfg.Nats.Address, hostName)
			if err != nil {
				zapLogger.Fatal("err connect nats", zap.Error(err))
			}
			return bus, nil
		}),
		// subsciptions
		fx.Provide(
			subscriptions.NewNotificationSubscription,
		),
		// invoke init func
		fx.Invoke(
			run,
			runWorker,
			func(r *gin.Engine, pe *prometheus.Exporter) {
				if pe != nil {
					r.GET("/metrics", func(c *gin.Context) {
						pe.ServeHTTP(c.Writer, c.Request)
					})
				}
			},
			// controller, register routes
			controller.RegisterAuthController,
			controller.RegisterGroupController,
			controller.RegisterRoleController,
			controller.RegisterUserController,
			controller.RegisterUploadController,
			controller.RegisterRuleController,
		),
	}

	err = fx.ValidateApp(opts...)
	if err != nil {
		log.Panic("err provide autowire", err)
	}

	if !cfg.Debug {
		opts = append(opts, fx.NopLogger)
	}

	fx.New(
		opts...,
	).Run()
}

func run(
	lifecycle fx.Lifecycle,
	zapLogger *zap.Logger,
	cfg *pkgConfig.BaseConfig,
	r *gin.Engine,
	p *profiler.Profiler,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					err := r.Run(cfg.Port)
					if err != nil {
						zapLogger.Fatal("r.Run", zap.String("port", cfg.Port), zap.Error(err))
					}
				}()

				zapLogger.Info("Hydra is running", zap.String("port", cfg.Port), zap.String("env", cfg.Env), zap.Bool("debug", cfg.Debug))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				zapLogger.Info("shutting down server...")

				if p != nil {
					err := p.Stop()
					if err != nil {
						zapLogger.Error("err stop profiler", zap.Error(err))
					}
				}

				return nil
			},
		},
	)
}

func runWorker(
	lifecycle fx.Lifecycle,
	zapLogger *zap.Logger,
	busFactory nats.BusFactory,
	notifySub *subscriptions.NotificationSubscription,
) {
	closeSubs := func() {}

	subscribe := func() error {
		zapLogger.Info("start worker, subscribe to nats server...")

		subs := []stan.Subscription{}

		notifySubs, err := notifySub.Subscribe()
		if err != nil {
			return fmt.Errorf("err notifySub.Subscribe: %w", err)
		}

		subs = append(subs, notifySubs...)

		zapLogger.Info("subscribe done", zap.Int("totalSub", len(subs)))
		closeSubs = func() {
			zapLogger.Info("close all subscriptions...")
			for _, sub := range subs {
				_ = sub.Close()
			}
		}

		return nil
	}

	err := subscribe()
	if err != nil {
		zapLogger.Panic("err subscribe worker", zap.Error(err))
	}

	busFactory.RegisterCallbackConnected([]nats.Callback{subscribe})

	lifecycle.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				closeSubs()
				return nil
			},
		},
	)
}
