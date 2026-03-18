package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"my-openclaw/pkg/types"

	"github.com/gorilla/websocket"
)

type Server struct {
	sessions *SessionManager
	agent    *AgentRuntime
	clients  map[*websocket.Conn]*types.AuthContext
	mu       sync.RWMutex
	seq      int64
	upgrader websocket.Upgrader
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*types.Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{sessions: make(map[string]*types.Session)}
}

func (m *SessionManager) GetOrCreate(key string) *types.Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s := m.sessions[key]; s != nil {
		s.LastActive = time.Now()
		return s
	}
	s := &types.Session{
		Key:        key,
		AgentID:    "main",
		Transcript: []types.Message{},
		Memory:     "",
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}
	m.sessions[key] = s
	return s
}

func (m *SessionManager) AddMessage(key string, msg types.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := m.sessions[key]
	if s == nil {
		s = &types.Session{Key: key, AgentID: "main", Transcript: []types.Message{}, CreatedAt: time.Now(), LastActive: time.Now()}
		m.sessions[key] = s
	}
	s.Transcript = append(s.Transcript, msg)
	s.LastActive = time.Now()
}

func (m *SessionManager) List() []*types.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*types.Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		result = append(result, s)
	}
	return result
}

type AgentRuntime struct {
	model string
}

func NewAgentRuntime(model string) *AgentRuntime {
	if model == "" {
		model = "mock"
	}
	return &AgentRuntime{model: model}
}

type Response struct {
	Content string
}

func (r *AgentRuntime) Run(sessionID string, messages []types.Message, memory string) Response {
	var userMsg string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			userMsg = messages[i].Content
			break
		}
	}
	return r.generateResponse(userMsg, memory)
}

func (r *AgentRuntime) generateResponse(userMsg, memory string) Response {
	userMsgL := ""
	if userMsg != "" {
		userMsgL = userMsg
	}

	content := ""
	switch {
	case contains(userMsgL, "hello") || contains(userMsgL, "hi"):
		content = "Hello! I am your AI assistant running on My-OpenClaw Gateway (Go + React). How can I help you today?"
	case contains(userMsgL, "who are you") || contains(userMsgL, "what are you"):
		content = "I am a minimal AI assistant built with Go backend and React frontend. Running model: " + r.model
	case contains(userMsgL, "remember"):
		content = "I remember! Your long-term memory: \"" + memory + "\""
	case contains(userMsgL, "help"):
		content = `I'm a minimal OpenClaw implementation!

Try saying:
- "hello" - Basic greeting
- "who are you" - About me
- "remember something" - Test memory
- "help" - Show this message

In full OpenClaw, I can:
- Use tools (browser, shell)
- Connect to multiple chat platforms
- Run proactively via heartbeat
- Use skills for capabilities`
	default:
		if userMsgL == "" {
			content = "Hello! How can I help you?"
		} else {
			content = "I received: \"" + userMsg + "\"\n\nThis is a Go + React implementation of OpenClaw Gateway. Try saying \"help\" to see what I can do!"
		}
	}
	return Response{Content: content}
}

func (r *AgentRuntime) GetModel() string {
	return r.model
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func New(model string) *Server {
	return &Server{
		sessions: NewSessionManager(),
		agent:    NewAgentRuntime(model),
		clients:  make(map[*websocket.Conn]*types.AuthContext),
		seq:      0,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ws":
		s.handleWebSocket(w, r)
	case "/health":
		s.handleHealth(w, r)
	case "/status":
		s.handleStatus(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"name":    "My-OpenClaw Gateway",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"ws":     "/ws",
				"health": "/health",
				"status": "/status",
			},
		})
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "version": "1.0.0"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	clientCount := len(s.clients)
	s.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"version":  "1.0.0",
		"sessions": len(s.sessions.List()),
		"clients":  clientCount,
		"model":    s.agent.GetModel(),
	})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	auth := types.AuthContext{OK: true, Role: "operator", Scopes: []string{"operator.read", "operator.write"}}
	s.mu.Lock()
	s.clients[conn] = &auth
	s.mu.Unlock()

	log.Printf("Client connected: %s", auth.Role)
	s.send(conn, types.EventFrame{Type: "event", Event: "gateway.ready", Payload: map[string]any{"version": "1.0.0"}, Seq: s.nextSeq()})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var frame types.RequestFrame
		if err := json.Unmarshal(msg, &frame); err != nil {
			s.sendError(conn, frame.ID, "Invalid JSON")
			continue
		}
		s.handleMessage(conn, &frame, auth)
	}

	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
	conn.Close()
}

func (s *Server) handleMessage(conn *websocket.Conn, frame *types.RequestFrame, auth types.AuthContext) {
	if frame.Type != "req" {
		s.sendError(conn, frame.ID, "Invalid frame type")
		return
	}

	var result any
	errMsg := ""

	switch frame.Method {
	case "connect":
		result = map[string]any{"type": "hello-ok", "protocol": 1, "policy": map[string]any{"tickIntervalMs": 15000}, "capabilities": []string{"chat", "sessions", "status"}}
	case "chat.send":
		result = s.handleChatSend(frame.Params)
	case "sessions.list":
		result = s.handleSessionsList()
	case "sessions.get":
		result = s.handleSessionsGet(frame.Params)
	case "status.get":
		result = map[string]any{"version": "1.0.0", "name": "My-OpenClaw Gateway", "model": s.agent.GetModel(), "sessions": len(s.sessions.List())}
	default:
		errMsg = "Unknown method: " + frame.Method
	}

	if errMsg != "" {
		s.sendError(conn, frame.ID, errMsg)
		return
	}
	s.send(conn, types.ResponseFrame{Type: "res", ID: frame.ID, OK: true, Payload: result})
}

func (s *Server) handleChatSend(params map[string]any) any {
	sessionKey := "main"
	if v, ok := params["session"].(string); ok {
		sessionKey = v
	}
	message := ""
	if v, ok := params["message"].(string); ok {
		message = v
	}

	s.sessions.AddMessage(sessionKey, types.Message{Role: "user", Content: message, Timestamp: time.Now().UnixMilli()})
	ses := s.sessions.GetOrCreate(sessionKey)
	resp := s.agent.Run(sessionKey, ses.Transcript, ses.Memory)
	s.sessions.AddMessage(sessionKey, types.Message{Role: "assistant", Content: resp.Content, Timestamp: time.Now().UnixMilli()})

	return map[string]any{"messageId": "msg_" + string(rune(time.Now().Unix())), "content": resp.Content, "session": sessionKey}
}

func (s *Server) handleSessionsList() any {
	list := s.sessions.List()
	sessions := make([]map[string]any, len(list))
	for i, ses := range list {
		sessions[i] = map[string]any{"key": ses.Key, "agentId": ses.AgentID, "messageCount": len(ses.Transcript), "lastActive": ses.LastActive}
	}
	return map[string]any{"sessions": sessions}
}

func (s *Server) handleSessionsGet(params map[string]any) any {
	key, _ := params["session"].(string)
	ses := s.sessions.GetOrCreate(key)
	return map[string]any{"key": ses.Key, "agentId": ses.AgentID, "messageCount": len(ses.Transcript), "memoryLength": len(ses.Memory)}
}

func (s *Server) send(conn *websocket.Conn, frame any) {
	if err := conn.WriteJSON(frame); err != nil {
		log.Printf("Send error: %v", err)
	}
}

func (s *Server) sendError(conn *websocket.Conn, id string, err string) {
	s.send(conn, types.ResponseFrame{Type: "res", ID: id, OK: false, Error: err})
}

func (s *Server) nextSeq() int64 {
	s.seq++
	return s.seq
}
