import _ from 'lodash';
import createError from 'http-errors';
import als from 'async-local-storage';

import commonErrors from './common';
import userErrors from './user';
import {app} from '../config';

const errors = {
  common: commonErrors,
  user: userErrors,
};

// 根据key生成自定义error
function get(key) {
  const lang = als.get('lang') || 'en';
  const item = _.get(errors, key);
  const err = new Error(item[lang] || 'Unknown error');
  return createError(item.status || 500, err, {
    code: `${app}-${item.code}`,
    expected: true,
  });
}

function create(...args) {
  const err = createError(...args);
  // 主动抛出的error设置expected，可以通过判断expected是否为true来识别是否为未知error
  err.expected = true;
  return err;
}

export default {
  get,
  create,
};
