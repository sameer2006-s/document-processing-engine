import { useCallback, useEffect, useState } from 'react'
import { listDocuments } from '../api/documents'
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

  const loadDocuments = useCallback(async (quiet = false) => {
    if (!quiet) {
      setError(null)
      setLoading(true)
    }
    try {
      const data = await listDocuments()
      setDocuments(data)
    } catch (err) {
      if (!quiet) setError(getErrorMessage(err))
    } finally {
      if (!quiet) setLoading(false)
    }
  }, [])

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

      <UploadZone onUploaded={() => void loadDocuments()} />

      {loading && <p className="status-text">Loading documents…</p>}
      {error && <div className="banner error">{error}</div>}

      {!loading && (
        <DocumentTable
          documents={documents}
          onDeleted={() => void loadDocuments()}
        />
      )}
    </div>
  )
}
