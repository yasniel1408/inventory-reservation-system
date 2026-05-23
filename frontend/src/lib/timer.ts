export function getRemainingSeconds(expiresAt: string, now: Date = new Date()): number {
  const diffMs = new Date(expiresAt).getTime() - now.getTime()
  return Math.max(0, Math.ceil(diffMs / 1000))
}

export function createExpiryNotifier(onExpire: () => void) {
  let notified = false

  return (remainingSeconds: number) => {
    if (remainingSeconds > 0 || notified) {
      return
    }
    notified = true
    onExpire()
  }
}
