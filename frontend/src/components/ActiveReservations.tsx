import { useEffect, useMemo, useState } from 'react'

import { createExpiryNotifier, getRemainingSeconds } from '../lib/timer'
import type { Reservation } from '../types/reservation'

type Props = {
  reservations: Reservation[]
  pendingReservationId: string | null
  onRelease: (reservation: Reservation) => void
  onExpired: () => void
}

export function ActiveReservations({ reservations, pendingReservationId, onRelease, onExpired }: Props) {
  if (reservations.length === 0) {
    return (
      <section className="panel" aria-labelledby="reservations-title">
        <div className="panel-heading">
          <h2 id="reservations-title">Reservas activas</h2>
          <span>0 activas</span>
        </div>
        <p className="empty-state">No hay reservas activas.</p>
      </section>
    )
  }

  return (
    <section className="panel" aria-labelledby="reservations-title">
      <div className="panel-heading">
        <h2 id="reservations-title">Reservas activas</h2>
        <span>{reservations.length} activas</span>
      </div>
      <div className="reservation-list">
        {reservations.map((reservation) => (
          <ReservationRow
            key={reservation.id}
            onExpired={onExpired}
            onRelease={onRelease}
            pending={pendingReservationId === reservation.id}
            reservation={reservation}
          />
        ))}
      </div>
    </section>
  )
}

function ReservationRow({ reservation, pending, onRelease, onExpired }: {
  reservation: Reservation
  pending: boolean
  onRelease: (reservation: Reservation) => void
  onExpired: () => void
}) {
  const [now, setNow] = useState(() => new Date())
  const notifyExpired = useMemo(() => createExpiryNotifier(onExpired), [onExpired])
  const remaining = getRemainingSeconds(reservation.expiresAt, now)

  useEffect(() => {
    notifyExpired(remaining)
  }, [notifyExpired, remaining])

  useEffect(() => {
    const id = window.setInterval(() => setNow(new Date()), 1000)
    return () => window.clearInterval(id)
  }, [])

  const label = reservation.itemName ?? reservation.itemId

  return (
    <article className="reservation-row">
      <div>
        <h3>{label}</h3>
        <p>Reserva activa · {reservation.quantity} unidades</p>
      </div>
      <div className="reservation-actions">
        <span aria-label={`Tiempo restante ${label}`}>{remaining}s</span>
        <button disabled={pending} onClick={() => onRelease(reservation)} type="button">
          {pending ? 'Liberando...' : `Liberar reserva ${label}`}
        </button>
      </div>
    </article>
  )
}
