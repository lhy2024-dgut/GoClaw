package channels

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goclaw/goclaw/ai"
	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（生产环境应限制）
	},
}

type WebChatChannel struct {
	logger   logs.Logger
	router   *gin.Engine
	aiEngine *ai.AIEngine
}

func NewWebChatChannel(logger logs.Logger, aiEngine *ai.AIEngine) *WebChatChannel {
	return &WebChatChannel{
		logger:   logger,
		aiEngine: aiEngine,
	}
}

func (w *WebChatChannel) Name() string {
	return "webchat"
}

func (w *WebChatChannel) Start(ctx context.Context, cfg *config.Config) error {
	w.logger.Info("WebChat channel starting...")

	// 注意：路由需要在HTTP服务器启动前注册
	// 这里我们只是标记通道已启动
	w.logger.Info("WebChat channel ready")

	// Simulate running
	go func() {
		<-ctx.Done()
		w.logger.Info("WebChat channel stopping")
	}()

	return nil
}

func (w *WebChatChannel) Stop() error {
	w.logger.Info("WebChat channel stopped")
	return nil
}

func (w *WebChatChannel) SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// WebSocket 端点
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			w.logger.Errorf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		w.logger.Info("WebSocket client connected")

		// 处理消息
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				w.logger.Errorf("WebSocket read error: %v", err)
				break
			}

			w.logger.Infof("Received message: %s", string(message))
			msgStr := string(message)

			// Get AI response (AI will handle function calling automatically)
			var response string
			if w.aiEngine != nil {
				response, err = w.aiEngine.Chat(context.Background(), msgStr)
				if err != nil {
					w.logger.Errorf("AI chat error: %v", err)
					response = "抱歉，处理您的请求时出错：" + err.Error()
				}
			} else {
				response = "AI 引擎未配置"
			}

			// Send response
			resp := map[string]interface{}{
				"message": response,
				"time":    time.Now().Format(time.RFC3339),
			}

			err = conn.WriteJSON(resp)
			if err != nil {
				w.logger.Errorf("WebSocket write error: %v", err)
				break
			}
		}
	})

	// WebChat 页面
	router.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"title": "GoClaw Chat",
		})
	})
}
