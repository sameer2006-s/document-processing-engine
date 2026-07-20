import { apiFetch } from './client'

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
