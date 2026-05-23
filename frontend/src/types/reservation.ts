export type InventoryItem = {
  id: string
  name: string
  totalStock: number
  reservedStock: number
  availableStock: number
}

export type ReservationStatus = 'active' | 'released' | 'expired' | 'confirmed'

export type Reservation = {
  id: string
  itemId: string
  itemName?: string
  quantity: number
  status: ReservationStatus
  expiresAt: string
  createdAt: string
  releasedAt?: string
}

export type ReleaseReservationResponse = {
  reservation: Pick<Reservation, 'id' | 'status'> & Partial<Reservation>
  stockReturned: boolean
  noop: boolean
}

export type ApiErrorBody = {
  error: {
    code: string
    message: string
    details: Record<string, unknown>
  }
}
