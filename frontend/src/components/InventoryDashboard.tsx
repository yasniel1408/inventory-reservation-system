import type { InventoryItem } from '../types/reservation'

type Props = {
  items: InventoryItem[]
  quantities: Record<string, number>
  pendingItemId: string | null
  onQuantityChange: (itemId: string, quantity: number) => void
  onReserve: (item: InventoryItem) => void
}

export function InventoryDashboard({ items, quantities, pendingItemId, onQuantityChange, onReserve }: Props) {
  if (items.length === 0) {
    return <p className="empty-state">No hay inventario disponible.</p>
  }

  return (
    <section className="panel inventory-panel" aria-labelledby="inventory-title">
      <div className="panel-heading">
        <h2 id="inventory-title">Inventario</h2>
        <span>{items.length} items</span>
      </div>
      <div className="inventory-grid">
        {items.map((item) => (
          <article className="inventory-row" key={item.id}>
            <div>
              <h3>{item.name}</h3>
              <dl className="stock-metrics">
                <div>
                  <dt>Total</dt>
                  <dd>{item.totalStock}</dd>
                </div>
                <div>
                  <dt>Reservado</dt>
                  <dd>{item.reservedStock}</dd>
                </div>
                <div>
                  <dt>Disponible</dt>
                  <dd>{item.availableStock}</dd>
                </div>
              </dl>
            </div>
            <div className="reserve-controls">
              <label htmlFor={`qty-${item.id}`}>Cantidad para {item.name}</label>
              <input
                id={`qty-${item.id}`}
                min="1"
                step="1"
                type="number"
                value={quantities[item.id] ?? 1}
                onChange={(event) => onQuantityChange(item.id, Number(event.target.value))}
              />
              <button
                disabled={pendingItemId === item.id || item.availableStock <= 0}
                onClick={() => onReserve(item)}
                type="button"
              >
                {pendingItemId === item.id ? 'Reservando...' : `Reservar ${item.name}`}
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  )
}
