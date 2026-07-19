import { useEffect, useState } from 'react'
import { setAuthToken, setOnUnauthorized } from './api/client'
import AuthPage from './components/AuthPage'
import Dashboard from './components/Dashboard'
import './App.css'

const TOKEN_KEY = 'token'

function App() {
  const [token, setToken] = useState<string | null>(() => {
    const stored = localStorage.getItem(TOKEN_KEY)
    setAuthToken(stored)
    return stored
  })

  useEffect(() => {
    setOnUnauthorized(() => {
      localStorage.removeItem(TOKEN_KEY)
      setAuthToken(null)
      setToken(null)
    })
  }, [])

  function handleAuth(nextToken: string) {
    localStorage.setItem(TOKEN_KEY, nextToken)
    setAuthToken(nextToken)
    setToken(nextToken)
  }

  function handleLogout() {
    localStorage.removeItem(TOKEN_KEY)
    setAuthToken(null)
    setToken(null)
  }

  return token ? (
    <Dashboard onLogout={handleLogout} />
  ) : (
    <AuthPage onAuth={handleAuth} />
  )
}

export default App
