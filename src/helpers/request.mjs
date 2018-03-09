import request from 'superagent';
import _ from 'lodash';

import httpStats from '../plugins/http-stats';

let timeout = 5 * 1000;

/**
 * Add default handle, timeout and httpStats
 *
 * @param {Request} req Superagent request
 */
function defaultHandle(req) {
  req.timeout(timeout);
  httpStats(req);
}

function setTimeout(v) {
  if (_.isNumber(v) && v >= 0) {
    timeout = v;
  }
}

/**
 * Superagent http get
 */
function get(...args) {
  const req = request.get(...args);
  defaultHandle(req);
  return req;
}

/**
 * Superagent http post
 */
function post(...args) {
  const req = request.post(...args);
  defaultHandle(req);
  return req;
}

/**
 * Superagent http put
 */
function put(...args) {
  const req = request.put(...args);
  defaultHandle(req);
  return req;
}

/**
 * Superagent http delete
 */
function del(...args) {
  const req = request.delete(...args);
  defaultHandle(req);
  return req;
}

/**
 * Superagent http patch
 */
function patch(...args) {
  const req = request.patch(...args);
  defaultHandle(req);
  return req;
}

export default {
  get,
  post,
  put,
  delete: del,
  patch,
  timeout: setTimeout,
};
