import _ from 'lodash';

/**
 * timeout中间件，可以配置超时时间与pass函数
 * @param  {Object} [options={}] options.timeout：超时间隔，单位ms；
 * options.pass：判断该请求是否需要跳过timeout，返回true则跳过
 * @return {Function} 返回中间件处理函数
 */
export default (options = {}) =>
  function timeout(ctx, next) {
    const pass = options.pass || _.noop;
    const ms = options.timeout || 5000;
    if (pass(ctx)) {
      return next();
    }
    let timer;
    return Promise.race([
      new Promise((resolve, reject) => {
        timer = setTimeout(() => {
          // 此异常设置为非主动抛出异常，方便定时排查超时接口
          const err = new Error('Request Timeout');
          err.status = 408;
          reject(err);
        }, ms);
      }),
      new Promise((resolve, reject) =>
        next().then(
          () => {
            clearTimeout(timer);
            resolve();
          },
          err => {
            clearTimeout(timer);
            reject(err);
          },
        ),
      ),
    ]);
  };
