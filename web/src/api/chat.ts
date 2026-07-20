import { apiFetch, getAuthToken, notifyUnauthorized, RequestError } from './client'

export interface ChatResponse {
  response: string
}

export function sendChat(documentId: string, userPrompt: string) {
  return apiFetch<ChatResponse>('/chat', {
    method: 'POST',
    body: JSON.stringify({
      document_id: documentId,
      user_prompt: userPrompt,
    }),
  })
}

type StreamHandlers = {
  onToken: (token: string) => void
  onDone?: () => void
  signal?: AbortSignal
}

/** Consume POST /chat/stream SSE: data {"token"}, event done / error. */
export async function streamChat(
  documentId: string,
  userPrompt: string,
  { onToken, onDone, signal }: StreamHandlers,
) {
  const headers = new Headers({
    'Content-Type': 'application/json',
    Accept: 'text/event-stream',
  })
  const token = getAuthToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  let response: Response
  try {
    response = await fetch('/chat/stream', {
      method: 'POST',
      headers,
      signal,
      body: JSON.stringify({
        document_id: documentId,
        user_prompt: userPrompt,
      }),
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw err
    }
    throw new RequestError(
      'Could not reach the API. Make sure the server is running on port 7070.',
      0,
    )
  }

  if (response.status === 401) {
    notifyUnauthorized()
    throw new RequestError('unauthorized', 401)
  }

  if (!response.ok) {
    let message = response.statusText
    try {
      const data = (await response.json()) as { error?: string }
      message = data.error ?? message
    } catch {
      // ignore
    }
    throw new RequestError(message, response.status)
  }

  if (!response.body) {
    throw new RequestError('No response body from stream', 500)
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let currentEvent = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const parts = buffer.split('\n\n')
    buffer = parts.pop() ?? ''

    for (const part of parts) {
      if (!part.trim()) continue

      currentEvent = ''
      let dataLine = ''
      for (const line of part.split('\n')) {
        if (line.startsWith('event:')) {
          currentEvent = line.slice(6).trim()
        } else if (line.startsWith('data:')) {
          dataLine = line.slice(5).trim()
        }
      }
      if (!dataLine) continue

      let payload: { token?: string; error?: string; ok?: boolean }
      try {
        payload = JSON.parse(dataLine) as {
          token?: string
          error?: string
          ok?: boolean
        }
      } catch {
        continue
      }

      if (currentEvent === 'error') {
        throw new RequestError(payload.error ?? 'stream error', 500)
      }
      if (currentEvent === 'done') {
        onDone?.()
        return
      }
      if (payload.token) {
        onToken(payload.token)
      }
    }
  }

  onDone?.()
}
