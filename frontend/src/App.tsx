import { useCallback, useEffect, useMemo, useState } from 'react'

import { createApiClient, ApiError, type ApiClient } from './api/client'
import { ActiveReservations } from './components/ActiveReservations'
import { InventoryDashboard } from './components/InventoryDashboard'
import type { InventoryItem, Reservation } from './types/reservation'
import './App.css'

type AppProps = {
  api?: ApiClient
}

type Feedback = {
  type: 'success' | 'error'
  message: string
} | null

function App({ api: injectedApi }: AppProps) {
  const api = useMemo(() => injectedApi ?? createApiClient(), [injectedApi])
  const [items, setItems] = useState<InventoryItem[]>([])
  const [reservations, setReservations] = useState<Reservation[]>([])
  const [quantities, setQuantities] = useState<Record<string, number>>({})
  const [loading, setLoading] = useState(true)
  const [feedback, setFeedback] = useState<Feedback>(null)
  const [pendingItemId, setPendingItemId] = useState<string | null>(null)
  const [pendingReservationId, setPendingReservationId] = useState<string | null>(null)

  const refresh = useCallback(async (showLoading = true) => {
    if (showLoading) {
      setLoading(true)
    }
    try {
      const [nextItems, nextReservations] = await Promise.all([
        api.listItems(),
        api.listReservations(),
      ])
      setItems(nextItems)
      setReservations(nextReservations)
    } catch {
      setFeedback({ type: 'error', message: 'No se pudo cargar el estado.' })
    } finally {
      setLoading(false)
    }
  }, [api])

  useEffect(() => {
    let isMounted = true

    async function loadInitialState() {
      try {
        const [nextItems, nextReservations] = await Promise.all([
          api.listItems(),
          api.listReservations(),
        ])
        if (!isMounted) {
          return
        }
        setItems(nextItems)
        setReservations(nextReservations)
      } catch {
        if (isMounted) {
          setFeedback({ type: 'error', message: 'No se pudo cargar el estado.' })
        }
      } finally {
        if (isMounted) {
          setLoading(false)
        }
      }
    }

    void loadInitialState()

    return () => {
      isMounted = false
    }
  }, [api])

  async function reserve(item: InventoryItem) {
    const quantity = quantities[item.id] ?? 1
    if (!Number.isInteger(quantity) || quantity <= 0) {
      setFeedback({ type: 'error', message: 'Ingresá una cantidad válida.' })
      return
    }

    setPendingItemId(item.id)
    setFeedback(null)
    try {
      await api.createReservation({
        itemId: item.id,
        quantity,
        idempotencyKey: crypto.randomUUID(),
      })
      setFeedback({ type: 'success', message: 'Reserva creada correctamente.' })
      await refresh()
    } catch (error) {
      setFeedback({ type: 'error', message: errorMessage(error) })
      await refresh()
    } finally {
      setPendingItemId(null)
    }
  }

  async function release(reservation: Reservation) {
    setPendingReservationId(reservation.id)
    setFeedback(null)
    try {
      await api.releaseReservation(reservation.id)
      setFeedback({ type: 'success', message: 'Reserva liberada.' })
      await refresh()
    } catch {
      setFeedback({ type: 'error', message: 'No se pudo liberar la reserva.' })
    } finally {
      setPendingReservationId(null)
    }
  }

  return (
    <main className="app-shell">
      <header className="topbar">
        <div>
          <p>Flash sale inventory</p>
          <h1>Reservas de inventario</h1>
        </div>
        <button onClick={() => void refresh()} type="button">Actualizar</button>
      </header>

      {feedback && (
        <div className={`feedback ${feedback.type}`} role="status">
          {feedback.message}
        </div>
      )}

      {loading ? (
        <p className="loading-state">Cargando inventario...</p>
      ) : (
        <div className="workspace">
          <InventoryDashboard
            items={items}
            onQuantityChange={(itemId, quantity) => setQuantities((current) => ({ ...current, [itemId]: quantity }))}
            onReserve={(item) => void reserve(item)}
            pendingItemId={pendingItemId}
            quantities={quantities}
          />
          <ActiveReservations
            onExpired={() => void refresh()}
            onRelease={(reservation) => void release(reservation)}
            pendingReservationId={pendingReservationId}
            reservations={reservations}
          />
        </div>
      )}
    </main>
  )
}

function errorMessage(error: unknown) {
  const code = error instanceof ApiError ? error.code : errorCode(error)
  if (code === 'insufficient_stock') {
    return 'No hay stock suficiente.'
  }
  if (code === 'invalid_quantity') {
    return 'Ingresá una cantidad válida.'
  }
  return 'No se pudo crear la reserva.'
}

function errorCode(error: unknown) {
  if (typeof error !== 'object' || error === null || !('code' in error)) {
    return null
  }
  const code = (error as { code: unknown }).code
  return typeof code === 'string' ? code : null
}

export default App
