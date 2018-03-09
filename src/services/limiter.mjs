/**
 * 此模块用于做次数限制的功能功能
 */

import Limiter from 'superlimiter';

import client from '../helpers/redis';

const loginFailLimiter = new Limiter(client, {
  // 间隔时间10分钟
  ttl: 10 * 60,
  max: 10,
  prefix: 'super-limiter-login-fail-',
});

/**
 * 获取登录失败数量
 *
 * @param {String} account
 * @returns {Number}
 */
export async function getLoginFailCount(account) {
  const count = await loginFailLimiter.getCount(account);
  return Number.parseInt(count || 0, 10);
}

/**
 * 登录失败数量增加
 *
 * @param {String} account
 */
export async function incLoginFailCount(account) {
  await loginFailLimiter.exec(account);
}

/**
 * 生成一个limiter
 *
 * @param {Object} options
 * @returns {Limiter}
 */
export function create(options) {
  return new Limiter(client, options);
}
