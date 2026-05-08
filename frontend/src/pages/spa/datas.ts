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
  { label: 'UTF-8', value: 'UTF-8' },
  // { label: 'ASCII', value: 'ASCII' },
  { label: 'HEX', value: 'HEX' },
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
  const encoder = new TextEncoder()
  const bytes = encoder.encode(data)
  return Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('')
}

export const hexToString = (hexString) => {
  const cleanHex = hexString.replace(/\s/g, '').toLowerCase()
  if (cleanHex.length % 2 !== 0) {
    throw new Error('无效的十六进制字符串：长度必须为偶数')
  }
  const bytes = new Uint8Array(cleanHex.length / 2)
  for (let i = 0; i < cleanHex.length; i += 2) {
    bytes[i / 2] = parseInt(cleanHex.substr(i, 2), 16)
  }
  const decoder = new TextDecoder('utf-8')
  return decoder.decode(bytes)
}