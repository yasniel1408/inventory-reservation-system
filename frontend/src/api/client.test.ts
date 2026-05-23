import { describe, expect, it, vi } from 'vitest'

import { ApiError, createApiClient } from './client'

describe('api client', () => {
  it('fetches inventory items', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      items: [{ id: 'item-1', name: 'Widget', totalStock: 10, reservedStock: 2, availableStock: 8 }],
    }))
    const api = createApiClient({ baseUrl: 'http://api.test', fetcher: fetchMock })

    const items = await api.listItems()

    expect(items).toHaveLength(1)
    expect(items[0].availableStock).toBe(8)
    expect(fetchMock).toHaveBeenCalledWith('http://api.test/items', expect.objectContaining({ method: 'GET' }))
  })

  it('sends idempotency key when creating a reservation', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      reservation: {
        id: 'res-1',
        itemId: 'item-1',
        quantity: 2,
        status: 'active',
        expiresAt: '2026-05-22T12:01:00Z',
        createdAt: '2026-05-22T12:00:00Z',
      },
    }, 201))
    const api = createApiClient({ baseUrl: 'http://api.test', fetcher: fetchMock })

    await api.createReservation({ itemId: 'item-1', quantity: 2, idempotencyKey: 'reserve-key' })

    expect(fetchMock).toHaveBeenCalledWith('http://api.test/reservations', expect.objectContaining({
      method: 'POST',
      headers: expect.objectContaining({ 'Idempotency-Key': 'reserve-key' }),
    }))
  })

  it('maps API errors to ApiError', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      error: { code: 'insufficient_stock', message: 'Not enough stock available.', details: {} },
    }, 409))
    const api = createApiClient({ baseUrl: 'http://api.test', fetcher: fetchMock })

    await expect(api.createReservation({ itemId: 'item-1', quantity: 99, idempotencyKey: 'reserve-key' }))
      .rejects.toMatchObject(new ApiError('insufficient_stock', 'Not enough stock available.', 409))
  })
})

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}
