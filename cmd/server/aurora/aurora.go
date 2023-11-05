package aurora

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nats-io/stan.go"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/controller"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/customer"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/subscriptions"
	pkgConfig "mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/jwtutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/logger"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/nats"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/rbacutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
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
		),
		// http server
		fx.Provide(transportutil.InitGinEngine),
		// route group
		fx.Provide(func(r *gin.Engine) *gin.RouterGroup {
			g := r.Group("aurora")
			g.GET("api-docs", routeutil.ServingDocs)

			return g
		}),
		// redis
		fx.Provide(func(cfg *configs.Config) (redis.Cmdable, error) {
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
		// monitoring
		fx.Provide(tracingutil.InitTelemetry),
		fx.Provide(func(cfg *configs.Config) (*profiler.Profiler, error) {
			if !cfg.RemoteProfiler.Enabled {
				return nil, nil
			}

			return profiler.Start(profiler.Config{
				ApplicationName: fmt.Sprintf("mm-printing.%s", cfg.Name),
				ServerAddress:   cfg.RemoteProfiler.ProfilerURL,
			})
		}),
		// repositories
		fx.Provide(
			repository.NewUserRepo,
			repository.NewPermissionRepo,
			repository2.NewCustomerRepo,
		),
		// services
		fx.Provide(

			func(zapLogger *zap.Logger, redisDB redis.Cmdable) ws.WebSocketService {
				hostName, _ := os.Hostname()
				return ws.NewApp(hostName, zapLogger, redisDB)
			},

			customer.NewService,
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
			subscriptions.NewZaloSubscription,
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
			controller.RegisterConfigController,
		),
	}

	if !cfg.Debug {
		opts = append(opts, fx.NopLogger)
	}

	err = fx.ValidateApp(opts...)
	if err != nil {
		log.Panic("err provide autowire", err)
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
	wsApp ws.WebSocketService,
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

				zapLogger.Info("Gezu is running", zap.String("port", cfg.Port), zap.String("env", cfg.Env), zap.Bool("debug", cfg.Debug))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				zapLogger.Info("shutting down server...")
				wsApp.HubStop()

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
	zaloEventSubscription *subscriptions.ZaloEventSubscription,
) {
	closeSubs := func() {}

	subscribe := func() error {
		zapLogger.Info("start worker, subscribe to nats server...")

		subs := []stan.Subscription{}

		zaloEventSubs, err := zaloEventSubscription.Subscribe()
		if err != nil {
			return fmt.Errorf("err zaloEventSubscription.Subscribe: %w", err)
		}

		subs = append(subs, zaloEventSubs...)

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
