export const parityModes = [
  { value: 0, label: 'None' },
  { value: 1, label: 'Odd' },
  { value: 2, label: 'Even' }
]

export const baudRates = [
  { value: 9600, label: 9600 },
  { value: 19200, label: 19200 },
  { value: 18400, label: 18400 }
]

export const dataBites = [
  { label: 5, value: 5 },
  { label: 6, value: 6 },
  { label: 7, value: 7 },
  { label: 8, value: 8 }
]

export const stopBites = [
  { label: 1, value: 1 },
  { label: 2, value: 2 }
]
export const codes = [
  { label: 'ASCII', value: 'ASCII' },
  { label: 'HEX', value: 'HEX' }
  // { label: 'UTF-8', value: 'UTF-8' },
  // { label: 'UTF-16', value: 'UTF-16' }
]


export function createPersistentTask(callback, interval, options: any = {}) {
  const { useBackgroundTime = false } = options
  let animationFrameId
  let lastExecTime = performance.now() // 记录上次执行的实际物理时间
  let isRunning = true
  function loop(currentTime) {
    if (!isRunning) return
    const deltaTime = currentTime - lastExecTime
    if (deltaTime >= interval) {
      callback(currentTime)
      lastExecTime += interval
      if (useBackgroundTime && deltaTime > interval * 2) {
        lastExecTime = currentTime
      }
    }
    animationFrameId = requestAnimationFrame(loop)
  }
  animationFrameId = requestAnimationFrame(loop)
  return () => {
    isRunning = false
    cancelAnimationFrame(animationFrameId)
  }
}

export const stringToHex = (data) => {
  if (typeof data === 'number') {
    const buffer = new ArrayBuffer(4)
    const view = new DataView(buffer)
    view.setUint32(0, data)
    return new Uint8Array(buffer)
  }
  if (typeof data === 'string') {
    return new TextEncoder().encode(data)
  }
  return new Uint8Array(0)
}

export const hexToString = (hex) => {
  const str = String(hex)
  return new TextEncoder().encode(str)
}