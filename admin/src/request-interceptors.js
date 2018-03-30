import request from 'axios';

import {urlPrefix} from './config';

request.interceptors.request.use(config => {
  if (!config.timeout) {
    config.timeout = 10 * 1000;
  }
  config.url = `${urlPrefix}${config.url}`;
  return config;
});
