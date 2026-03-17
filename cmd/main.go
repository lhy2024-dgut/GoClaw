package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goclaw",
	Short: "GoClaw - Personal AI Assistant",
	Long:  `GoClaw is a personal AI assistant that runs 24/7, connects to messaging apps, and can execute tasks.`,
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the GoClaw server",
	Long:  `Start the GoClaw server and run it as a daemon.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Initialize logger
		logger := logs.NewLogger(cfg.Log.Level)
		logger.Info("Starting GoClaw server...")

		// Create context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-sigChan
			logger.Infof("Received signal: %v", sig)
			cancel()
		}()

		// Start the server
		if err := startServer(ctx, cfg, logger); err != nil {
			logger.Errorf("Server error: %v", err)
			os.Exit(1)
		}

		logger.Info("GoClaw server stopped")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoClaw v0.1.0")
		fmt.Println("Personal AI Assistant")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
