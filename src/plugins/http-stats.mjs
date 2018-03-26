import _ from 'lodash';

import stringify from '../helpers/stringify';

/**
 * HTTP Request stats
 *
 * @param {Request} req Superagent request
 */
export default function httpStats(req) {
  const stats = {};
  const finished = _.once(() => {
    // get the use time of request
    stats.use = Date.now() - stats.startedAt;
    delete stats.startedAt;
    if (stats.error) {
      console.error(stringify(stats));
    } else {
      console.info(stringify(stats));
    }
    req.emit('stats', stats);
  });
  req.once('request', () => {
    // eslint-disable-next-line
    const sendData = req._data;
    Object.assign(stats, {
      host: req.host,
      path: req.req.path,
      method: req.method,
      startedAt: Date.now(),
    });
    const {backendServer} = req;
    if (backendServer) {
      Object.assign(stats, _.pick(backendServer, ['ip', 'port']));
    }
    if (sendData && !_.isEmpty(sendData)) {
      stats.data = stringify(sendData);
    }
  });
  req.once('error', err => {
    stats.code = -1;
    stats.error = err.message;
    finished();
  });
  req.once('response', res => {
    stats.code = res.statusCode;
    finished();
  });
}
