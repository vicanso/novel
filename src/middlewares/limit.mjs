/**
 * 对系统连接数做限制配置的中间件
 */

import limit from 'koa-connection-limit';

import * as limiterService from '../services/limiter';
import * as globals from '../helpers/globals';

/**
 * 设置不同的connecting数量级别，不同的连接数对应不同的状态。
 * 使用koa-connection-limit，根据当前连接数，当状态变化时，触发回调函数。
 * 当系统状态达到`high`时，设置系统`status`为`pause`。当连接数降低，不再是`high`时，
 * 延时`interval`将系统重置为`running`
 * @param  {Object} options {mid: Integer, high: Integer}
 * @param  {Integer} interval 重置延时间隔
 * @return {Function} 返回中间件处理函数
 * @see {@link https://github.com/vicanso/koa-connection-limit|GitHub}
 */
export const connection = (options, interval) => {
  let connectionLimitTimer;
  return limit(options, status => {
    console.info(`connection-limit status:${status}`);
    if (status === 'high') {
      // 如果并发处理数已到达high值，设置状态为 pause，此时ping请求返回error，反向代理(nginx, haproxy)认为此程序有问题，不再转发请求到此程序
      globals.pause();
      /* istanbul ignore if */
      if (connectionLimitTimer) {
        clearTimeout(connectionLimitTimer);
        connectionLimitTimer = null;
      }
    } else if (!globals.isRunning()) {
      // 状态为low或者mid时，延时interval ms将服务设置回running
      connectionLimitTimer = setTimeout(() => {
        globals.start();
        connectionLimitTimer = null;
      }, interval);
      connectionLimitTimer.unref();
    }
  });
};

/**
 * 创建一个limiter中间件
 */
export const createLimiter = options => {
  const limiter = limiterService.create(options);
  return limiter.middleware();
};
