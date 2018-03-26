import als from 'async-local-storage';
import bluebird from 'bluebird';

import {addReplacer} from './helpers/stringify';
import './helpers/joi-extend';
import './helpers/logger';

// 堆的全局的Promise对象
global.Promise = bluebird;

// set stringify mask
addReplacer(key => {
  const reg = /password/gi;
  if (reg.test(key)) {
    return '***';
  }
  return null;
});

// 启用async local storage
als.enable();
