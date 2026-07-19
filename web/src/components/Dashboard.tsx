import { useCallback, useEffect, useState } from 'react'
import { listDocuments, searchMyFiles } from '../api/documents'
import { getErrorMessage } from '../api/client'
import type { Document } from '../api/types'
import DocumentTable from './DocumentTable'
import UploadZone from './UploadZone'

interface DashboardProps {
  onLogout: () => void
}

const IN_FLIGHT = new Set(['pending', 'ocr-processing'])

export default function Dashboard({ onLogout }: DashboardProps) {
  const [documents, setDocuments] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchInput, setSearchInput] = useState('')
  const [query, setQuery] = useState('')

  const loadDocuments = useCallback(async (quiet = false, searchQuery = query) => {
    if (!quiet) {
      setError(null)
      setLoading(true)
    }
    try {
      const trimmed = searchQuery.trim()
      const data = trimmed
        ? await searchMyFiles(trimmed)
        : await listDocuments()
      setDocuments(data)
    } catch (err) {
      if (!quiet) setError(getErrorMessage(err))
    } finally {
      if (!quiet) setLoading(false)
    }
  }, [query])

  useEffect(() => {
    const id = window.setTimeout(() => {
      setQuery(searchInput.trim())
    }, 300)
    return () => window.clearTimeout(id)
  }, [searchInput])

  useEffect(() => {
    void loadDocuments()
  }, [loadDocuments])

  // Poll while any document is still processing so status stays fresh.
  useEffect(() => {
    const hasInFlight = documents.some((d) => IN_FLIGHT.has(d.status))
    if (!hasInFlight) return

    const id = window.setInterval(() => {
      void loadDocuments(true)
    }, 2000)

    return () => window.clearInterval(id)
  }, [documents, loadDocuments])

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div>
          <h1>Documents</h1>
          <p>Upload, download, and manage your files.</p>
        </div>
        <button type="button" className="btn ghost" onClick={onLogout}>
          Sign out
        </button>
      </header>

      <UploadZone
        onUploaded={() => {
          setSearchInput('')
          setQuery('')
          void loadDocuments(false, '')
        }}
      />

      <div className="search-bar">
        <label className="field search-field">
          <span>Search</span>
          <input
            type="search"
            placeholder="Search by file name or OCR text…"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            autoComplete="off"
          />
        </label>
        {searchInput && (
          <button
            type="button"
            className="btn ghost"
            onClick={() => setSearchInput('')}
          >
            Clear
          </button>
        )}
      </div>

      {query && !loading && documents.length > 0 && (
        <p className="status-text">
          {documents.length} result{documents.length === 1 ? '' : 's'} for “
          {query}”
        </p>
      )}

      {loading && <p className="status-text">Loading documents…</p>}
      {error && <div className="banner error">{error}</div>}

      {!loading && (
        <DocumentTable
          documents={documents}
          onDeleted={() => void loadDocuments()}
          emptyMessage={
            query
              ? `No files matched “${query}”.`
              : 'No documents yet. Upload your first file above.'
          }
        />
      )}
    </div>
  )
}
