import request from 'superagent';
import Koa from 'koa';
import Stream from 'stream';
import _ from 'lodash';
import ms from 'ms';

import * as config from '../config';

/**
 * 从参数列表中获取第一个符合条件的参数返回，如果都不符合，则使用默认值返回
 * @param  {Array} arr 参数列表
 * @param  {Function} validate 校验函数，如果符合则返回true
 * @param  {Any} defaultValue 默认值
 * @return {Any} 返回符合条件的值或者默认值
 * @example
 * const utils = require('./helpers/utils');
 * // max: 100
 * const max = utils.getParam(['name', true, 100], _.isInteger, 10);
 */
export function getParam(arr, validate, defaultValue) {
  const v = arr.find(validate);
  if (v === undefined) {
    return defaultValue;
  }
  return v;
}

/**
 * 请求当前应用
 * @param {String} method http request method
 * @param {String} url http request url
 */
export function selfRequest(method, url) {
  const requestUrl = `http://127.0.0.1:${config.port}${url}`;
  return request[method](requestUrl);
}

/**
 * 延时执行(Promise)
 * @param {Integer} ms 延时的ms
 */
export function delay(v) {
  return new Promise(resolve => setTimeout(resolve, v));
}

/**
 * 该请求是否不可以缓存
 * @param {Context} ctx
 */
export function isNoCache(ctx) {
  const {method} = ctx;
  if (method !== 'GET' && method !== 'HEAD') {
    return true;
  }
  if (
    ctx.get('Cache-Control') === 'no-cache' ||
    ctx.query['cache-control'] === 'no-cache'
  ) {
    return true;
  }
  return false;
}

export function isDevelopment() {
  return config.env === 'development';
}

export function isProduction() {
  return config.env === 'production';
}

export function setCache(ctx, ttl, sMaxAge) {
  let seconds = ttl;
  if (_.isString(seconds)) {
    seconds = _.ceil(ms(ttl) / 1000);
  }
  let cacheControl = `public, max-age=${seconds}`;
  if (sMaxAge) {
    let sMaxAgeSeconds = sMaxAge;
    if (_.isString(sMaxAgeSeconds)) {
      sMaxAgeSeconds = _.ceil(ms(sMaxAgeSeconds) / 1000);
    }
    cacheControl += `, s-maxage=${sMaxAgeSeconds}`;
  }
  ctx.set('Cache-Control', cacheControl);
}

export function createContext(requestUrl) {
  const socket = new Stream.Duplex();
  const req = Object.assign({headers: {}, socket}, Stream.Readable.prototype);
  const res = Object.assign({_headers: {}, socket}, Stream.Writable.prototype);
  req.socket.remoteAddress = req.socket.remoteAddress || '127.0.0.1';
  const app = new Koa();
  // eslint-disable-next-line
  res.getHeader = k => res._headers[k.toLowerCase()];
  // eslint-disable-next-line
  res.setHeader = (k, v) => (res._headers[k.toLowerCase()] = v);
  // eslint-disable-next-line
  res.removeHeader = k => delete res._headers[k.toLowerCase()];
  const ctx = app.createContext(req, res);
  ctx.params = {};
  if (requestUrl) {
    ctx.request.url = requestUrl;
  }
  ctx.setCache = (ttl, sMaxAge) => {
    setCache(ctx, ttl, sMaxAge);
  };
  return ctx;
}
