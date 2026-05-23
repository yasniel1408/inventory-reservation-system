import { describe, expect, it, vi } from 'vitest'

import { getRemainingSeconds, createExpiryNotifier } from './timer'

describe('timer logic', () => {
  it('calculates remaining seconds without going negative', () => {
    const now = new Date('2026-05-22T12:00:00Z')

    expect(getRemainingSeconds('2026-05-22T12:00:30Z', now)).toBe(30)
    expect(getRemainingSeconds('2026-05-22T11:59:59Z', now)).toBe(0)
  })

  it('notifies expiration once', () => {
    const onExpire = vi.fn()
    const notifier = createExpiryNotifier(onExpire)

    notifier(1)
    notifier(0)
    notifier(0)

    expect(onExpire).toHaveBeenCalledTimes(1)
  })
})
