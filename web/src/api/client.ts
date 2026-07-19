import type { ApiError } from './types'

let authToken: string | null = null
let onUnauthorized: (() => void) | null = null

export function setAuthToken(token: string | null) {
  authToken = token
}

export function setOnUnauthorized(callback: () => void) {
  onUnauthorized = callback
}

export class RequestError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'RequestError'
    this.status = status
  }
}

async function readErrorMessage(response: Response) {
  let message = response.statusText
  try {
    const data = (await response.json()) as ApiError
    message = data.error ?? message
  } catch {
    // ignore non-JSON error bodies
  }
  return message
}

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const headers = new Headers(options.headers)

  if (authToken) {
    headers.set('Authorization', `Bearer ${authToken}`)
  }

  if (options.body && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  let response: Response
  try {
    response = await fetch(path, { ...options, headers })
  } catch {
    throw new RequestError(
      'Could not reach the API. Make sure the server is running on port 7070.',
      0,
    )
  }

  if (response.status === 401) {
    const message = await readErrorMessage(response)
    if (authToken) {
      onUnauthorized?.()
    }
    throw new RequestError(message, 401)
  }

  if (!response.ok) {
    const message = await readErrorMessage(response)
    throw new RequestError(message, response.status)
  }

  if (response.status === 204) {
    return undefined as T
  }

  const contentType = response.headers.get('Content-Type') ?? ''
  if (contentType.includes('application/json')) {
    try {
      return (await response.json()) as T
    } catch {
      throw new RequestError('Received an invalid response from the server.', 200)
    }
  }

  return response.blob() as Promise<T>
}

export function getErrorMessage(err: unknown) {
  if (err instanceof RequestError) {
    return err.message
  }
  if (err instanceof Error) {
    return err.message
  }
  return 'Something went wrong'
}
