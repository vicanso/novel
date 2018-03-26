import _ from 'lodash';

const replacers = [];
/**
 * 格式化字符串
 * @param data
 */
export default function stringify(data) {
  return JSON.stringify(data, (key, value) => {
    let result;
    _.forEach(replacers, replacer => {
      if (!_.isNil(result)) {
        return;
      }
      result = replacer(key, value);
    });
    if (_.isNil(result)) {
      return value;
    }
    return result;
  });
}

/**
 * 增加 replacer
 * @param fn
 */
export function addReplacer(fn) {
  replacers.push(fn);
}
