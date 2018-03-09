import als from 'async-local-storage';
import stringify from 'simple-stringify';
import bluebird from 'bluebird';

import './helpers/joi-extend';
import './helpers/logger';

// 堆的全局的Promise对象
global.Promise = bluebird;

// set stringify mask
stringify.isSecret = key => {
  const reg = /password/gi;
  return reg.test(key);
};
stringify.addFormat('_id', v => {
  if (!v) {
    return '';
  }
  return v.toString();
});

// 启用async local storage
als.enable();
