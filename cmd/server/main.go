package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"mmlabel.gitlab.com/mm-printing-backend/cmd/server/aurora"
	"mmlabel.gitlab.com/mm-printing-backend/cmd/server/hydra"
	"mmlabel.gitlab.com/mm-printing-backend/cmd/server/iot"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/logger"
	"mmlabel.gitlab.com/mm-printing-backend/version"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configPath  string
		migrateMode bool
	)
	cmdMigrate := &cobra.Command{
		Use:   "migrate",
		Short: "migrate database",
		Long:  `migrate database for all service`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			MigrateDB(configPath)
		},
	}
	cmdAurora := &cobra.Command{
		Use:   "aurora",
		Short: "Start aurora server",
		Long:  `Start aurora server to handle chat, message...`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			aurora.Run(ctx, configPath)
		},
	}
	cmdHydra := &cobra.Command{
		Use:   "hydra",
		Short: "Start hydra server",
		Long:  `Start hydra server to handle authenticates, user profile...`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			hydra.Run(ctx, configPath)
		},
	}
	cmdIot := &cobra.Command{
		Use:   "iot",
		Short: "Start iot handler",
		Long:  `Start iot handler to handle iot device`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			iot.Run(ctx, configPath)
		},
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of mm-printing backend",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("mm-printing backend " + version.Version)
			fmt.Println("Git hash " + version.GitHash)
			fmt.Println("Go version " + version.GoVersion)
		},
	}

	rootCmd := &cobra.Command{Use: "mm-printing"}
	rootCmd.AddCommand(
		cmdMigrate,
		cmdHydra,
		cmdAurora,
		cmdIot,
		cmdVersion,
	)

	rootCmd.PersistentFlags().StringVar(
		&configPath,
		"configPath",
		"",
		"path to configuration file, usually used for configuration",
	)

	rootCmd.PersistentFlags().BoolVar(
		&migrateMode,
		"migrateMode",
		false,
		"run migrateMode",
	)

	err := rootCmd.MarkPersistentFlagRequired("configPath")
	if err != nil {
		log.Panic(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Panic(err)
	}
}

func MigrateDB(configPath string) {
	cfg, err := configs.NewConfig(configPath)
	if err != nil {
		log.Panic(err)
	}

	uri := string(cfg.CockroachDB.URI)
	uri = strings.ReplaceAll(uri, "postgres://", "crdb-postgres://")
	m, err := migrate.New(
		cfg.CockroachDB.MigrationPath,
		uri,
	)
	if err != nil {
		log.Panic("err init migrate", err)
	}

	zapLogger := logger.NewZapLogger(cfg.Logger.LogLevel, cfg.Env == "local")
	m.Log = &logger.WrapLog{
		L: zapLogger,
	}
	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Panic("err migrate up: ", err)
		}
		log.Println("no migration file change")
	}
}
