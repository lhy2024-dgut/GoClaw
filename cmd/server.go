package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goclaw/goclaw/ai"
	"github.com/goclaw/goclaw/channels"
	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
	"github.com/goclaw/goclaw/tools"
)

func startServer(ctx context.Context, cfg *config.Config, logger logs.Logger) error {
	// Initialize tool registry
	toolRegistry := tools.NewRegistry(logger)
	toolRegistry.Register(tools.NewFileSystemTool())

	// Initialize AI engine
	aiEngine := ai.NewAIEngine(logger, cfg, toolRegistry)

	// Initialize channels manager
	channelManager := channels.NewManager(logger)

	// Start all configured channels
	if err := channelManager.StartAll(ctx, cfg); err != nil {
		return fmt.Errorf("failed to start channels: %w", err)
	}

	// Set AI engine for WebChat channel
	channelManager.SetAIEngine(aiEngine)

	// Start HTTP server
	httpServer := startHTTPServer(cfg, logger, channelManager, aiEngine)

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
	}

	// Stop all channels
	if err := channelManager.StopAll(); err != nil {
		logger.Errorf("Failed to stop channels: %v", err)
	}

	return nil
}

func startHTTPServer(cfg *config.Config, logger logs.Logger, channelManager *channels.Manager, aiEngine *ai.AIEngine) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Load templates
	router.LoadHTMLGlob("templates/*")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Metrics endpoint (basic)
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        cfg.App.Name,
			"version":     cfg.App.Version,
			"description": "Personal AI Assistant",
			"endpoints": []string{
				"GET  /",
				"GET  /health",
				"GET  /metrics",
				"POST /api/chat",
				"GET  /api/chat/history",
			},
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.POST("/chat", func(c *gin.Context) {
			var request struct {
				Message string `json:"message"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Get AI response
			aiResponse, err := aiEngine.Chat(c.Request.Context(), request.Message)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": aiResponse,
			})
		})

		api.GET("/chat/history", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Chat history coming soon",
			})
		})
	}

	// Setup WebChat routes
	webchatChan := channelManager.GetChannel("webchat")
	if webchatChan != nil {
		if webchat, ok := webchatChan.(*channels.WebChatChannel); ok {
			webchat.SetupRoutes(router, cfg)
		}
	}

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Infof("HTTP server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP server error: %v", err)
		}
	}()

	return srv
}
