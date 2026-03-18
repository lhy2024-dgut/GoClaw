import { useState, useEffect, useRef } from 'react'

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
}

interface GatewayStatus {
  connected: boolean
  version: string
  sessions: number
}

function App() {
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [status, setStatus] = useState<GatewayStatus>({ connected: false, version: '1.0.0', sessions: 0 })
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    connectWebSocket()
    return () => {
      ws?.close()
    }
  }, [])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const connectWebSocket = () => {
    const socket = new WebSocket('ws://127.0.0.1:18789/ws')

    socket.onopen = () => {
      console.log('Connected to Gateway')
      setStatus(s => ({ ...s, connected: true }))
      
      // Send connect request
      socket.send(JSON.stringify({
        type: 'req',
        id: 'init',
        method: 'connect',
        params: { client: { id: 'web', version: '1.0.0', platform: 'web' } }
      }))
    }

    socket.onmessage = (event) => {
      try {
        const frame = JSON.parse(event.data)
        
        if (frame.type === 'res' && frame.ok && frame.payload?.type === 'hello-ok') {
          console.log('Handshake complete')
        }
        
        if (frame.type === 'res' && frame.id === 'chat') {
          setLoading(false)
          if (frame.ok && frame.payload) {
            const msg: Message = {
              id: frame.payload.messageId || Date.now().toString(),
              role: 'assistant',
              content: frame.payload.content || ''
            }
            setMessages(prev => [...prev, msg])
          } else {
            console.error('Chat error:', frame.error)
          }
        }
      } catch (e) {
        console.error('Parse error:', e)
      }
    }

    socket.onclose = () => {
      setStatus(s => ({ ...s, connected: false }))
      console.log('Disconnected, reconnecting...')
      setTimeout(connectWebSocket, 3000)
    }

    socket.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    setWs(socket)
  }

  const sendMessage = async () => {
    if (!input.trim() || !ws || loading) return

    const userMsg: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input.trim()
    }
    setMessages(prev => [...prev, userMsg])
    setInput('')
    setLoading(true)

    ws.send(JSON.stringify({
      type: 'req',
      id: 'chat',
      method: 'chat.send',
      params: {
        session: 'main',
        message: userMsg.content
      }
    }))
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  return (
    <div className="app">
      <header className="header">
        <h1>🦞 My-OpenClaw</h1>
        <div className="status">
          <span className={`indicator ${status.connected ? 'connected' : 'disconnected'}`}></span>
          {status.connected ? 'Connected' : 'Disconnected'}
        </div>
      </header>

      <main className="chat">
        {messages.length === 0 && (
          <div className="welcome">
            <p>👋 Hello! I'm your AI assistant.</p>
            <p>Try saying "hello", "help", or ask me anything!</p>
          </div>
        )}
        
        {messages.map(msg => (
          <div key={msg.id} className={`message ${msg.role}`}>
            <div className="message-content">{msg.content}</div>
          </div>
        ))}
        
        {loading && (
          <div className="message assistant">
            <div className="message-content typing">Thinking...</div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </main>

      <footer className="input-area">
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type a message..."
          disabled={!status.connected || loading}
        />
        <button onClick={sendMessage} disabled={!status.connected || loading || !input.trim()}>
          Send
        </button>
      </footer>
    </div>
  )
}

export default App
