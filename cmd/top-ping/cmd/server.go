package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"top-ping/internal/app/router"
	"top-ping/pkg/database"
	"top-ping/pkg/logger"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start http server",
	Run: func(cmd *cobra.Command, args []string) {
		// Create context that listens for the interrupt signal from the OS.
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		defer logger.Sync()
		profile := config.GetString("application.profile")

		var logging logger.Config
		loggingErr := config.UnmarshalKey("logging", &logging)
		if loggingErr != nil {
			panic("loading logging configuration error!!!")
		}
		logger.Init(profile, &logging)

		var mysqlConf database.DatasourceConfig
		databaseErr := config.UnmarshalKey("mysql", &mysqlConf)
		if databaseErr != nil {
			panic("loading logging configuration error!!!")
		}
		database.Init(&mysqlConf)

		host := config.GetString("server.host")
		port := config.GetInt("server.port")
		addr := fmt.Sprintf("%s:%d", host, port)

		logger.Infof(ctx, "Server: listening on: %s", addr)
		srv := &http.Server{
			Addr:         addr,
			Handler:      router.Router(profile, &logging),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  30 * time.Second,
		}

		// Initializing the server in a goroutine so that
		// it won't block the graceful shutdown handling below
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalf(ctx, "listen: %s", err)
			}
		}()

		// Listen for the interrupt signal.
		<-ctx.Done()

		// Restore default behavior on the interrupt signal and notify user of shutdown.
		stop()
		//booted.Logger.Println("shutting down gracefully, press Ctrl+C again to force")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf(ctx, "Server forced to shutdown: %v", err)
		}
	},
}
