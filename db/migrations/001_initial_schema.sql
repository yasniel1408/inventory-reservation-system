BEGIN;

CREATE TABLE IF NOT EXISTS items (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  total_stock INTEGER NOT NULL,
  reserved_stock INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT items_total_stock_nonnegative CHECK (total_stock >= 0),
  CONSTRAINT items_reserved_stock_nonnegative CHECK (reserved_stock >= 0),
  CONSTRAINT items_reserved_stock_lte_total_stock CHECK (reserved_stock <= total_stock)
);

CREATE TABLE IF NOT EXISTS reservations (
  id UUID PRIMARY KEY,
  item_id UUID NOT NULL REFERENCES items(id) ON DELETE RESTRICT,
  session_id TEXT NOT NULL,
  quantity INTEGER NOT NULL,
  status TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  released_at TIMESTAMPTZ,
  expired_at TIMESTAMPTZ,
  CONSTRAINT reservations_quantity_positive CHECK (quantity > 0),
  CONSTRAINT reservations_status_allowed CHECK (status IN ('active', 'released', 'expired', 'confirmed')),
  CONSTRAINT reservations_expires_at_after_creation CHECK (expires_at >= created_at),
  CONSTRAINT reservations_terminal_timestamps_consistent CHECK (
    (status = 'active' AND released_at IS NULL AND expired_at IS NULL)
    OR (status = 'released' AND released_at IS NOT NULL AND expired_at IS NULL)
    OR (status = 'expired' AND expired_at IS NOT NULL AND released_at IS NULL)
    OR (status = 'confirmed' AND released_at IS NULL AND expired_at IS NULL)
  )
);

CREATE TABLE IF NOT EXISTS idempotency_keys (
  scope TEXT NOT NULL,
  key TEXT NOT NULL,
  request_hash TEXT NOT NULL,
  status TEXT NOT NULL,
  response_code INTEGER,
  response_body JSONB,
  reservation_id UUID REFERENCES reservations(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT idempotency_keys_status_allowed CHECK (status IN ('in_progress', 'completed', 'failed')),
  CONSTRAINT idempotency_keys_response_code_valid CHECK (
    response_code IS NULL OR (response_code >= 100 AND response_code <= 599)
  ),
  CONSTRAINT idempotency_keys_response_consistent CHECK (
    (
      status = 'in_progress'
      AND response_code IS NULL
      AND response_body IS NULL
      AND reservation_id IS NULL
    )
    OR (
      status IN ('completed', 'failed')
      AND response_code IS NOT NULL
      AND response_body IS NOT NULL
    )
  ),
  CONSTRAINT idempotency_keys_scope_key_unique UNIQUE (scope, key)
);

CREATE INDEX IF NOT EXISTS idx_reservations_active_expires_at
  ON reservations (expires_at)
  WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_reservations_session_status_created_at
  ON reservations (session_id, status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_reservations_item_status
  ON reservations (item_id, status);

CREATE INDEX IF NOT EXISTS idx_idempotency_keys_reservation_id
  ON idempotency_keys (reservation_id);

COMMIT;
