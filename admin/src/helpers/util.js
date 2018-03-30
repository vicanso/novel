import _ from 'lodash';

import {sha256} from './crypto';
import {app} from '../config';

// 获取出错信息
export function getErrorMessage(err) {
  let message = err;
  if (err && err.response) {
    const {data, headers} = err.response;
    const id = headers['x-response-id'];
    // eslint-disable-next-line
    const code = data.code.replace(`${app}-`, '');
    message = `${data.message}(${code}) [${id}]`;
  }
  if (_.isError(message)) {
    message = message.message;
  }
  if (err.code === 'ECONNABORTED') {
    message = '请求超时，请重新再试';
  }
  return message;
}

// 生成密码
export function genPassword(account, password) {
  const pwd = sha256(password);
  return sha256(`${account}-${pwd}-${app}`);
}

// 获取日期格式化字符串
export function getDate(str) {
  const date = new Date(str);
  const fill = v => {
    if (v >= 10) {
      return `${v}`;
    }
    return `0${v}`;
  };
  const month = fill(date.getMonth() + 1);
  const day = fill(date.getDate());
  const hours = fill(date.getHours());
  const mintues = fill(date.getMinutes());
  const seconds = fill(date.getSeconds());
  return `${date.getFullYear()}-${month}-${day} ${hours}:${mintues}:${seconds}`;
}

// 等待ttl时长
export function waitFor(ttl, startedAt) {
  let delay = ttl;
  if (startedAt) {
    delay = ttl - (Date.now() - startedAt);
  }
  return new Promise(resolve => {
    setTimeout(resolve, Math.max(0, delay));
  });
}

// 获取不一致的数据
export function diff(original, current, keys) {
  const diffKeys = keys || _.keys(current);
  const result = {};
  _.forEach(diffKeys, key => {
    const value = current[key];
    if (original[key] !== value) {
      result[key] = value;
    }
  });
  return result;
}
