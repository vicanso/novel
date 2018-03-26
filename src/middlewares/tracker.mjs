/**
 * 用于跟踪用户行为的中间件，主要用于一些会修改数据的操作
 */
import _ from 'lodash';

import stringify from '../helpers/stringify';
import influx from '../helpers/influx';

/**
 * 记录用户的行为日志到influxdb中
 * @param  {Object} data 用户行为日志数据
 */
function logUserTracker(data) {
  console.info(`user tracker ${stringify(data)}`);
  const tags = 'category result'.split(' ');
  influx.write('userTracker', _.omit(data, tags), _.pick(data, tags));
}

/**
 * 生成行为日志中间件，根据设置的参数列表获取用户提交的参数，
 * 以后最后的结果，记录到influxdb中
 * @param  {String} category 该用户行为分类，如：用户注册、用户收藏
 * @return {Function} 返回中间件处理函数
 */
export default category =>
  async function userTracker(ctx, next) {
    const data = {
      category,
      ip: ctx.ip,
    };
    const account = _.get(ctx, 'session.user.account');
    if (account) {
      data.account = account;
    }
    const params = _.extend({}, ctx.query, ctx.request.body, ctx.params);
    if (!_.isEmpty(params)) {
      data.params = stringify(params).replace(/"/g, '\\"');
    }
    const start = Date.now();
    const resultLog = (use, result) => {
      data.result = result;
      data.use = use;
      logUserTracker(data);
    };
    let type = 'fail';
    try {
      await next();
      if (ctx.status !== 404) {
        type = 'success';
      }
    } catch (err) {
      data.message = err.message;
      throw err;
    } finally {
      resultLog(Date.now() - start, type);
    }
  };
