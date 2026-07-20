import { useEffect, useId, useRef, useState } from 'react'
import { deleteDocument, downloadFile, fetchThumbnail } from '../api/documents'
import { getErrorMessage } from '../api/client'
import type { Document } from '../api/types'

interface DocumentTableProps {
  documents: Document[]
  onDeleted: () => void
  onChat?: (doc: Document) => void
  emptyMessage?: string
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
    case 'thumbnail-processing':
      return 'Thumbnail'
    case 'tag-processing':
      return 'Tagging'
    case 'done':
    case 'ocr-done':
    case 'thumbnail-done':
    case 'tag-done':
      return 'Done'
    case 'failed':
      return 'Failed'
    default:
      return status
  }
}

/** Backend may store raw OCR text or JSON {"text":"..."}. */
function formatOcrResult(raw?: string) {
  if (!raw?.trim()) return null
  try {
    const parsed = JSON.parse(raw) as { text?: unknown }
    if (typeof parsed?.text === 'string') {
      return parsed.text.trim() || null
    }
  } catch {
    // plain text
  }
  return raw.trim()
}

function DocumentThumbnail({ doc }: { doc: Document }) {
  const [src, setSrc] = useState<string | null>(null)

  useEffect(() => {
    if (!doc.thumbnail_key) {
      setSrc(null)
      return
    }

    let objectUrl: string | null = null
    let cancelled = false

    void (async () => {
      try {
        const blob = await fetchThumbnail(doc.id)
        if (cancelled) return
        objectUrl = URL.createObjectURL(blob)
        setSrc(objectUrl)
      } catch {
        if (!cancelled) setSrc(null)
      }
    })()

    return () => {
      cancelled = true
      if (objectUrl) URL.revokeObjectURL(objectUrl)
    }
  }, [doc.id, doc.thumbnail_key])

  if (!doc.thumbnail_key || !src) {
    return <div className="doc-thumb placeholder" aria-hidden />
  }

  return (
    <img
      className="doc-thumb"
      src={src}
      alt=""
      width={48}
      height={48}
    />
  )
}

function OcrViewer({
  title,
  text,
  onClose,
}: {
  title: string
  text: string
  onClose: () => void
}) {
  const titleId = useId()
  const closeRef = useRef<HTMLButtonElement>(null)
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    closeRef.current?.focus()
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'

    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => {
      document.body.style.overflow = prev
      window.removeEventListener('keydown', onKey)
    }
  }, [onClose])

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1500)
    } catch {
      // ignore clipboard failures
    }
  }

  return (
    <div
      className="ocr-modal-backdrop"
      role="presentation"
      onClick={onClose}
    >
      <div
        className="ocr-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        onClick={(e) => e.stopPropagation()}
      >
        <header className="ocr-modal-header">
          <div>
            <h2 id={titleId}>OCR result</h2>
            <p className="ocr-modal-subtitle">{title}</p>
          </div>
          <div className="ocr-modal-actions">
            <button type="button" className="btn ghost" onClick={() => void handleCopy()}>
              {copied ? 'Copied' : 'Copy'}
            </button>
            <button
              ref={closeRef}
              type="button"
              className="btn ghost"
              onClick={onClose}
            >
              Close
            </button>
          </div>
        </header>
        <div className="ocr-modal-body">
          <pre className="ocr-modal-text">{text}</pre>
        </div>
      </div>
    </div>
  )
}

export default function DocumentTable({
  documents,
  onDeleted,
  onChat,
  emptyMessage = 'No documents yet. Upload your first file above.',
}: DocumentTableProps) {
  const [busyId, setBusyId] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [viewer, setViewer] = useState<{ title: string; text: string } | null>(
    null,
  )

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
        <p>{emptyMessage}</p>
      </div>
    )
  }

  return (
    <div className="table-wrap">
      {error && <div className="banner error">{error}</div>}
      <table className="doc-table">
        <thead>
          <tr>
            <th aria-label="Preview" />
            <th>Name</th>
            <th>Status</th>
            <th>Tags</th>
            <th>OCR</th>
            <th>Size</th>
            <th>Uploaded</th>
            <th aria-label="Actions" />
          </tr>
        </thead>
        <tbody>
          {documents.map((doc) => {
            const ocrText = formatOcrResult(doc.ocr_result)
            const tags = doc.tags?.filter(Boolean) ?? []

            return (
              <tr key={doc.id}>
                <td className="thumb-cell">
                  <DocumentThumbnail doc={doc} />
                </td>
                <td className="name-cell">{doc.original_name}</td>
                <td>
                  <span className={`status-badge status-${doc.status}`}>
                    {statusLabel(doc.status)}
                  </span>
                </td>
                <td className="tags-cell">
                  {tags.length > 0 ? (
                    <div className="tag-list">
                      {tags.map((tag) => (
                        <span key={tag} className="tag-chip">
                          {tag}
                        </span>
                      ))}
                    </div>
                  ) : (
                    <span className="ocr-empty">—</span>
                  )}
                </td>
                <td className="ocr-cell">
                  {ocrText ? (
                    <button
                      type="button"
                      className="btn ghost ocr-view-btn"
                      onClick={() =>
                        setViewer({ title: doc.original_name, text: ocrText })
                      }
                    >
                      View text
                    </button>
                  ) : (
                    <span className="ocr-empty">—</span>
                  )}
                </td>
                <td>{formatSize(doc.file_size)}</td>
                <td>{formatDate(doc.created_at)}</td>
                <td className="actions-cell">
                  {onChat && (
                    <button
                      type="button"
                      className="btn ghost"
                      disabled={busyId === doc.id || !ocrText}
                      title={
                        ocrText
                          ? 'Chat about this document'
                          : 'OCR text required before chat'
                      }
                      onClick={() => onChat(doc)}
                    >
                      Chat
                    </button>
                  )}
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
            )
          })}
        </tbody>
      </table>

      {viewer && (
        <OcrViewer
          title={viewer.title}
          text={viewer.text}
          onClose={() => setViewer(null)}
        />
      )}
    </div>
  )
}
