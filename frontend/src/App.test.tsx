import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'

import App from './App'
import type { ApiClient } from './api/client'

describe('reservation app', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('renders inventory and creates a reservation with feedback', async () => {
    const api = fakeApi()
    render(<App api={api} />)

    expect(await screen.findByText('Standard Widget')).toBeInTheDocument()
    await userEvent.clear(screen.getByLabelText('Cantidad para Standard Widget'))
    await userEvent.type(screen.getByLabelText('Cantidad para Standard Widget'), '2')
    await userEvent.click(screen.getByRole('button', { name: 'Reservar Standard Widget' }))

    expect(await screen.findByText('Reserva creada correctamente.')).toBeInTheDocument()
    expect(api.createReservation).toHaveBeenCalledWith(expect.objectContaining({
      itemId: 'item-1',
      quantity: 2,
    }))
    expect(api.listItems).toHaveBeenCalledTimes(2)
    expect(api.listReservations).toHaveBeenCalledTimes(2)
  })

  it('shows insufficient stock error without adding a fake reservation', async () => {
    const api = fakeApi()
    api.createReservation = vi.fn().mockRejectedValue({ code: 'insufficient_stock', message: 'Not enough stock available.' })
    render(<App api={api} />)

    expect(await screen.findByText('Standard Widget')).toBeInTheDocument()
    await userEvent.clear(screen.getByLabelText('Cantidad para Standard Widget'))
    await userEvent.type(screen.getByLabelText('Cantidad para Standard Widget'), '99')
    await userEvent.click(screen.getByRole('button', { name: 'Reservar Standard Widget' }))

    expect(await screen.findByText('No hay stock suficiente.')).toBeInTheDocument()
    expect(screen.queryByText('Reserva activa')).not.toBeInTheDocument()
  })

  it('releases active reservations and refetches state', async () => {
    const api = fakeApi({
      reservations: [{
        id: 'res-1',
        itemId: 'item-1',
        itemName: 'Standard Widget',
        quantity: 1,
        status: 'active',
        expiresAt: '2026-05-22T12:01:00Z',
        createdAt: '2026-05-22T12:00:00Z',
      }],
    })
    render(<App api={api} />)

    expect(await screen.findByRole('button', { name: 'Liberar reserva Standard Widget' })).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Liberar reserva Standard Widget' }))

    expect(await screen.findByText('Reserva liberada.')).toBeInTheDocument()
    expect(api.releaseReservation).toHaveBeenCalledWith('res-1')
  })
})

function fakeApi(overrides: Partial<{ reservations: Awaited<ReturnType<ApiClient['listReservations']>> }> = {}): ApiClient {
  return {
    listItems: vi.fn().mockResolvedValue([
      { id: 'item-1', name: 'Standard Widget', totalStock: 10, reservedStock: 0, availableStock: 10 },
    ]),
    createReservation: vi.fn().mockResolvedValue({
      id: 'res-1',
      itemId: 'item-1',
      quantity: 2,
      status: 'active',
      expiresAt: '2026-05-22T12:01:00Z',
      createdAt: '2026-05-22T12:00:00Z',
    }),
    listReservations: vi.fn().mockResolvedValue(overrides.reservations ?? []),
    releaseReservation: vi.fn().mockResolvedValue({
      reservation: { id: 'res-1', status: 'released' },
      stockReturned: true,
      noop: false,
    }),
  }
}
