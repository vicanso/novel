const map = new Map();

/**
 * Judge the application is running
 * @returns {Boolean}
 */
export function isRunning() {
  return map.get('status') === 'running';
}

/**
 * Pause the application
 */
export function pause() {
  map.set('status', 'pause');
}

/**
 * Start the application
 */
export function start() {
  map.set('status', 'running');
}

/**
 * Get connecting count
 */
export function getConnectingCount() {
  return map.get('connectingCount') || 0;
}

/**
 * Set connecting count
 */
export function setConnectingCount(count) {
  let concurrency = 'low';
  if (count > 1000) {
    concurrency = 'high';
  } else if (count > 500) {
    concurrency = 'mid';
  }
  map.set('concurrency', concurrency);
  map.set('connectingCount', count);
}

/**
 * Get the concurrency
 */
export function getConcurrency() {
  return map.get('concurrency') || 'low';
}

/**
 * Set the performance
 */
export function setPerformance(value) {
  map.set('performance', value);
}

/**
 * Get the performance
 */
export function getPerformance() {
  return map.get('performance');
}
