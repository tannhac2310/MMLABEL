package iot

import (
	"context"
	"fmt"
	"log"
	"os"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/device_config"
	"mmlabel.gitlab.com/mm-printing-backend/internal/iot/subscriptions"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/option"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/product_quality"

	"github.com/go-redis/redis"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/department"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/device"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_export"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_import"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_return"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_order"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_order_stage_device"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/stage"

	repository2 "mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	pkgConfig "mmlabel.gitlab.com/mm-printing-backend/pkg/configs"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/customer"
	"mmlabel.gitlab.com/mm-printing-backend/internal/iot/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/logger"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/nats"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/tracingutil"
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
			repository2.NewProductionOrderRepo,
			repository2.NewProductionOrderStageRepo,
			repository2.NewStageRepo,
			repository2.NewDeviceRepo,
			repository2.NewDepartmentRepo,
			repository2.NewProductionOrderStageDeviceRepo,
			repository2.NewCustomFieldRepo,
			repository2.NewInkRepo,
			repository2.NewInkImportRepo,
			repository2.NewInkImportDetailRepo,
			repository2.NewInkExportRepo,
			repository2.NewInkExportDetailRepo,
			repository2.NewInkReturnRepo,
			repository2.NewInkReturnDetailRepo,
			repository2.NewHistoryRepo,
			//repository2.NewProductQualityRepo,
			repository2.NewOptionRepo,
			repository2.NewProductionOrderDeviceConfigRepo,
			repository2.NewDeviceProgressStatusHistoryRepo,
			repository2.NewDeviceBrokenHistoryRepo,
			repository2.NewDeviceWorkingHistoryRepo,
			repository.NewUserRoleRepo,
			repository.NewRoleRepo,
			repository.NewRolePermissionRepo,
			repository2.NewProductionOrderStageResponsibleRepo,
		),
		// services
		fx.Provide(

			func(zapLogger *zap.Logger, redisDB redis.Cmdable) ws.WebSocketService {
				hostName, _ := os.Hostname()
				return ws.NewApp(hostName, zapLogger, redisDB)
			},
			customer.NewService,
			production_order.NewService,
			stage.NewService,
			device.NewService,
			department.NewService,
			production_order_stage_device.NewService,
			ink.NewService,
			ink_import.NewService,
			ink_export.NewService,
			ink_return.NewService,
			product_quality.NewService,
			option.NewService,
			device_config.NewService,
			role.NewService,
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
			subscriptions.NewMQTTSubscription,
		),
		// invoke init func
		fx.Invoke(
			run,
			runWorker,

			// controller, register routes
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
	wsApp ws.WebSocketService,
	p *profiler.Profiler,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					//err := r.Run(cfg.Port)
					//if err != nil {
					//	zapLogger.Fatal("r.Run", zap.String("port", cfg.Port), zap.Error(err))
					//}
				}()

				zapLogger.Info("IOT is running", zap.String("port", cfg.Port), zap.String("env", cfg.Env), zap.Bool("debug", cfg.Debug))
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
	eventMQTTSubscription *subscriptions.EventMQTTSubscription,
) {
	closeSubs := func() {}

	subscribe := func() error {
		fmt.Println("==============================================MQTT==============================================")
		zapLogger.Info("start worker, subscribe to nats server...")

		err := eventMQTTSubscription.Subscribe()
		if err != nil {
			return fmt.Errorf("err eventMQTTSubscription.Subscribe: %w", err)
		}

		zapLogger.Info("subscribe done")

		closeSubs = func() {
			zapLogger.Info("close all subscriptions...")

		}

		return nil
	}

	err := subscribe()
	if err != nil {
		zapLogger.Panic("err subscribe worker", zap.Error(err))
	}

	lifecycle.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				closeSubs()
				return nil
			},
		},
	)
}
