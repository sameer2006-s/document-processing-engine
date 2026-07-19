import { useState } from 'react'
import { deleteDocument, downloadFile } from '../api/documents'
import { getErrorMessage } from '../api/client'
import type { Document } from '../api/types'

interface DocumentTableProps {
  documents: Document[]
  onDeleted: () => void
}

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function formatDate(value: string) {
  return new Date(value).toLocaleString()
}

function statusLabel(status: string) {
  switch (status) {
    case 'pending':
      return 'Pending'
    case 'ocr-processing':
      return 'Processing'
    case 'done':
      return 'Done'
    case 'failed':
      return 'Failed'
    default:
      return status
  }
}

export default function DocumentTable({
  documents,
  onDeleted,
}: DocumentTableProps) {
  const [busyId, setBusyId] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  async function handleDownload(doc: Document) {
    setError(null)
    setBusyId(doc.id)
    try {
      const blob = await downloadFile(doc.id)
      const url = URL.createObjectURL(blob)
      const anchor = document.createElement('a')
      anchor.href = url
      anchor.download = doc.original_name
      anchor.click()
      URL.revokeObjectURL(url)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setBusyId(null)
    }
  }

  async function handleDelete(doc: Document) {
    if (!window.confirm(`Delete "${doc.original_name}"?`)) return

    setError(null)
    setBusyId(doc.id)
    try {
      await deleteDocument(doc.id)
      onDeleted()
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setBusyId(null)
    }
  }

  if (documents.length === 0) {
    return (
      <div className="empty-state">
        <p>No documents yet. Upload your first file above.</p>
      </div>
    )
  }

  return (
    <div className="table-wrap">
      {error && <div className="banner error">{error}</div>}
      <table className="doc-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Status</th>
            <th>Size</th>
            <th>Uploaded</th>
            <th aria-label="Actions" />
          </tr>
        </thead>
        <tbody>
          {documents.map((doc) => (
            <tr key={doc.id}>
              <td className="name-cell">{doc.original_name}</td>
              <td>
                <span className={`status-badge status-${doc.status}`}>
                  {statusLabel(doc.status)}
                </span>
              </td>
              <td>{formatSize(doc.file_size)}</td>
              <td>{formatDate(doc.created_at)}</td>
              <td className="actions-cell">
                <button
                  type="button"
                  className="btn ghost"
                  disabled={busyId === doc.id}
                  onClick={() => void handleDownload(doc)}
                >
                  Download
                </button>
                <button
                  type="button"
                  className="btn danger"
                  disabled={busyId === doc.id}
                  onClick={() => void handleDelete(doc)}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
