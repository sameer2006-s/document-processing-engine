import { useEffect, useRef, useState, type FormEvent, type KeyboardEvent } from 'react'
import { streamChat } from '../api/chat'
import { getErrorMessage } from '../api/client'
import type { Document } from '../api/types'

interface ChatPageProps {
  doc: Document
  onBack: () => void
}

interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

export default function ChatPage({ doc, onBack }: ChatPageProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [prompt, setPrompt] = useState('')
  const [sending, setSending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)
  const abortRef = useRef<AbortController | null>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, sending])

  useEffect(() => {
    inputRef.current?.focus()
    return () => {
      abortRef.current?.abort()
    }
  }, [])

  async function handleSend(e: FormEvent) {
    e.preventDefault()
    const text = prompt.trim()
    if (!text || sending) return

    setError(null)
    setPrompt('')
    setMessages((prev) => [
      ...prev,
      { role: 'user', content: text },
      { role: 'assistant', content: '' },
    ])
    setSending(true)

    const controller = new AbortController()
    abortRef.current = controller

    try {
      await streamChat(doc.id, text, {
        signal: controller.signal,
        onToken: (token) => {
          setMessages((prev) => {
            const next = [...prev]
            const last = next[next.length - 1]
            if (last?.role === 'assistant') {
              next[next.length - 1] = {
                ...last,
                content: last.content + token,
              }
            }
            return next
          })
        },
      })
    } catch (err) {
      if (err instanceof DOMException && err.name === 'AbortError') {
        return
      }
      setError(getErrorMessage(err))
      setMessages((prev) => {
        const next = [...prev]
        const last = next[next.length - 1]
        if (last?.role === 'assistant' && last.content === '') {
          next.pop()
          if (next[next.length - 1]?.role === 'user') {
            next.pop()
          }
        }
        return next
      })
      setPrompt(text)
    } finally {
      abortRef.current = null
      setSending(false)
    }
  }

  function handleKeyDown(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      void handleSend(e)
    }
  }

  return (
    <div className="chat-page">
      <header className="chat-header">
        <button type="button" className="btn ghost" onClick={onBack}>
          ← Documents
        </button>
        <div className="chat-header-meta">
          <h1>Chat</h1>
          <p>{doc.original_name}</p>
        </div>
      </header>

      <div className="chat-thread" role="log" aria-live="polite">
        {messages.length === 0 && !sending && (
          <div className="chat-empty">
            <p>
              Ask a question about this document. Answers use its OCR text as
              context.
            </p>
          </div>
        )}
        {messages.map((msg, i) => (
          <div key={`${msg.role}-${i}`} className={`chat-bubble ${msg.role}`}>
            <span className="chat-role">
              {msg.role === 'user' ? 'You' : 'Assistant'}
            </span>
            <pre className="chat-content">
              {msg.content ||
                (sending && i === messages.length - 1 ? '…' : '')}
            </pre>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      {error && <div className="banner error">{error}</div>}

      <form className="chat-composer" onSubmit={(e) => void handleSend(e)}>
        <textarea
          ref={inputRef}
          rows={3}
          placeholder="Ask about this document…"
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={sending}
        />
        <button
          type="submit"
          className="btn primary"
          disabled={sending || !prompt.trim()}
        >
          {sending ? 'Streaming…' : 'Send'}
        </button>
      </form>
    </div>
  )
}
