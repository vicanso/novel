import Redis from 'ioredis';
import _ from 'lodash';

import * as config from '../config';

const client = new Redis(config.redisUri, {
  keyPrefix: `${config.app}:`,
});

const delayLog = _.throttle((message, type) => {
  const maskUri = config.redisUri.replace(/:\S+@/, '//:***@');
  if (type === 'error') {
    console.alert(`${maskUri} error, ${message})`);
  } else {
    console.info(`${maskUri} ${message}`);
  }
}, 3000);

client.on('error', err => delayLog(err.message, 'error'));

// 延时输出日志，避免一直断开连接时大量无用日志
client.on('connect', () => delayLog('connected'));

class SessionStore {
  constructor(redisClient) {
    this.redisClient = redisClient;
  }
  async get(key) {
    const data = await this.redisClient.get(key);
    if (!data) {
      return null;
    }
    return JSON.parse(data);
  }
  async set(key, json, maxAge) {
    await this.redisClient.psetex(key, maxAge, JSON.stringify(json));
  }
  async destroy(key) {
    await this.redisClient.del(key);
  }
}

export default client;

export const sessionStore = new SessionStore(client);

// 获取key并锁定ttl时长
export async function lock(key, ttl) {
  if (!key || !ttl) {
    throw new Error('key and ttl can not be null');
  }
  const result = await client.set(`${key}-lock`, 1, 'NX', 'EX', ttl);
  if (result === 'OK') {
    return true;
  }
  return false;
}
