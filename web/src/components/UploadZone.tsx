import { useRef, useState, type DragEvent } from 'react'
import { uploadFile } from '../api/documents'
import { getErrorMessage } from '../api/client'

interface UploadZoneProps {
  onUploaded: () => void
}

export default function UploadZone({ onUploaded }: UploadZoneProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [message, setMessage] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  async function handleFiles(files: FileList | null) {
    const file = files?.[0]
    if (!file) return

    setUploading(true)
    setMessage(null)
    setError(null)

    try {
      await uploadFile(file)
      setMessage(`Uploaded ${file.name}`)
      onUploaded()
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setUploading(false)
    }
  }

  function onDrop(event: DragEvent) {
    event.preventDefault()
    setDragging(false)
    void handleFiles(event.dataTransfer.files)
  }

  return (
    <div className="upload-section">
      <div
        className={dragging ? 'upload-zone dragging' : 'upload-zone'}
        onDragOver={(e) => {
          e.preventDefault()
          setDragging(true)
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={onDrop}
        onClick={() => inputRef.current?.click()}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault()
            inputRef.current?.click()
          }
        }}
      >
        <input
          ref={inputRef}
          type="file"
          hidden
          onChange={(e) => void handleFiles(e.target.files)}
        />
        <p className="upload-title">
          {uploading ? 'Uploading…' : 'Drop a file here or click to browse'}
        </p>
        <p className="upload-hint">Any file type supported</p>
      </div>

      {message && <div className="banner success">{message}</div>}
      {error && <div className="banner error">{error}</div>}
    </div>
  )
}
