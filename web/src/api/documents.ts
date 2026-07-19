import { apiFetch } from './client'
import type { Document } from './types'

export function listDocuments() {
  return apiFetch<Document[]>('/documents')
}

export function uploadFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return apiFetch<Document>('/upload', {
    method: 'POST',
    body: formData,
  })
}

export function downloadFile(id: string) {
  return apiFetch<Blob>(`/get-file/${id}`)
}

export function deleteDocument(id: string) {
  return apiFetch<{ message: string }>(`/documents/${id}`, {
    method: 'DELETE',
  })
}
