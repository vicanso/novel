/**
 * 用于对返回数据输出至日志中
 * 此中间件能方便查询接口的响应数据，但也会有可能输出了不希望输出的数据，需要慎用
 */

import _ from 'lodash';
import stringify from 'simple-stringify';

export default (level = 2) => async (ctx, next) => {
  await next();
  if (ctx.state.logResponse && ctx.body && _.isObject(ctx.body)) {
    const logLevel = _.get(ctx, 'state.logResponseLevel', level);
    console.info(`response: ${stringify.json(ctx.body, logLevel)}`);
  }
};
