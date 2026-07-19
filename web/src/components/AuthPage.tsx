import { useState, type FormEvent } from 'react'
import { login, register } from '../api/auth'
import { getErrorMessage } from '../api/client'

type AuthMode = 'signin' | 'register'

interface AuthPageProps {
  onAuth: (token: string) => void
}

export default function AuthPage({ onAuth }: AuthPageProps) {
  const [mode, setMode] = useState<AuthMode>('signin')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setError(null)
    setLoading(true)

    try {
      if (mode === 'signin') {
        const response = await login({ email, password })
        onAuth(response.token)
      } else {
        await register({ email, password, firstName, lastName })
        const response = await login({ email, password })
        onAuth(response.token)
      }
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-brand">
          <img src="/favicon.svg" alt="" className="auth-logo" />
          <h1>Documents</h1>
          <p>Sign in to upload and manage your files.</p>
        </div>

        <div className="tabs" role="tablist">
          <button
            type="button"
            role="tab"
            className={mode === 'signin' ? 'tab active' : 'tab'}
            aria-selected={mode === 'signin'}
            onClick={() => setMode('signin')}
          >
            Sign in
          </button>
          <button
            type="button"
            role="tab"
            className={mode === 'register' ? 'tab active' : 'tab'}
            aria-selected={mode === 'register'}
            onClick={() => setMode('register')}
          >
            Create account
          </button>
        </div>

        {error && <div className="banner error">{error}</div>}

        <form className="auth-form" onSubmit={handleSubmit}>
          {mode === 'register' && (
            <div className="field-row">
              <label className="field">
                <span>First name</span>
                <input
                  type="text"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  required
                  autoComplete="given-name"
                />
              </label>
              <label className="field">
                <span>Last name</span>
                <input
                  type="text"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  required
                  autoComplete="family-name"
                />
              </label>
            </div>
          )}

          <label className="field">
            <span>Email</span>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
            />
          </label>

          <label className="field">
            <span>Password</span>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete={
                mode === 'signin' ? 'current-password' : 'new-password'
              }
            />
          </label>

          <button type="submit" className="btn primary" disabled={loading}>
            {loading
              ? 'Please wait…'
              : mode === 'signin'
                ? 'Sign in'
                : 'Create account'}
          </button>
        </form>
      </div>
    </div>
  )
}
