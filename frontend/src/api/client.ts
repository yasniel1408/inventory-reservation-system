import type { ApiErrorBody, InventoryItem, ReleaseReservationResponse, Reservation } from '../types/reservation'

export type CreateReservationInput = {
  itemId: string
  quantity: number
  idempotencyKey: string
}

export type ApiClient = {
  listItems: () => Promise<InventoryItem[]>
  createReservation: (input: CreateReservationInput) => Promise<Reservation>
  listReservations: () => Promise<Reservation[]>
  releaseReservation: (id: string) => Promise<ReleaseReservationResponse>
}

type ClientOptions = {
  baseUrl?: string
  fetcher?: typeof fetch
}

export class ApiError extends Error {
  code: string
  status: number

  constructor(code: string, message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

const defaultBaseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export function createApiClient(options: ClientOptions = {}): ApiClient {
  const baseUrl = (options.baseUrl ?? defaultBaseUrl).replace(/\/$/, '')
  const fetcher = options.fetcher ?? fetch

  return {
    async listItems() {
      const body = await request<{ items: InventoryItem[] }>(fetcher, `${baseUrl}/items`, { method: 'GET' })
      return body.items
    },
    async createReservation(input) {
      const body = await request<{ reservation: Reservation }>(fetcher, `${baseUrl}/reservations`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Idempotency-Key': input.idempotencyKey,
        },
        body: JSON.stringify({ itemId: input.itemId, quantity: input.quantity }),
      })
      return body.reservation
    },
    async listReservations() {
      const body = await request<{ reservations: Reservation[] }>(fetcher, `${baseUrl}/reservations`, { method: 'GET' })
      return body.reservations
    },
    async releaseReservation(id) {
      return request<ReleaseReservationResponse>(fetcher, `${baseUrl}/reservations/${id}`, { method: 'DELETE' })
    },
  }
}

async function request<T>(fetcher: typeof fetch, url: string, init: RequestInit): Promise<T> {
  const response = await fetcher(url, init)
  const body = await parseJson(response)

  if (!response.ok) {
    const errorBody = body as ApiErrorBody
    throw new ApiError(
      errorBody.error?.code ?? 'internal_error',
      errorBody.error?.message ?? 'Unexpected API error.',
      response.status,
    )
  }

  return body as T
}

async function parseJson(response: Response): Promise<unknown> {
  const text = await response.text()
  if (!text) {
    return {}
  }
  return JSON.parse(text)
}
