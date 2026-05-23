BEGIN;

INSERT INTO items (id, name, total_stock, reserved_stock)
VALUES
  ('00000000-0000-0000-0000-000000000001', 'Standard Widget', 100, 0),
  ('00000000-0000-0000-0000-000000000002', 'Limited Edition Widget', 2, 0)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  total_stock = EXCLUDED.total_stock,
  reserved_stock = EXCLUDED.reserved_stock,
  updated_at = NOW();

COMMIT;
