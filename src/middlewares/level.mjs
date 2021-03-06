/**
 * 用于生成根据当前系统运行级别而处理的中间件
 */
import _ from 'lodash';

import errors from '../errors';
import * as settingServices from '../services/setting';

/**
 * 对于某些非重要接口，可以设置不同的系统级别以上才做响应，如果系统运行级别低于设置级别，则直接返回出错。
 * 用于当系统负载较高的时候，系统自动或者手工将系统级别调低，降低系统负载
 * @param {Integer} level 当前中间件运行的系统级别
 * @return {Function} 返回中间件处理函数
 */
export default level => (ctx, next) => {
  const systemLevel = settingServices.get('limitLevel.level');
  if (!_.isNil(systemLevel) && systemLevel < level) {
    const err = errors.get('common.systemLevelLimit');
    err.message = err.message
      .replace('#{systemLevel}', systemLevel)
      .replace('#{level}', level);
    throw err;
  }
  return next();
};
