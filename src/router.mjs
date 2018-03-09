import _ from 'lodash';
import router from 'koa-router-parser';

import debug from './helpers/debug';
import middlewares from './middlewares';
import routes from './routes';
import controllers from './controllers';

function getRouter(descList) {
  return router.parse(descList);
}

function addToRouter(category, fns) {
  if (_.isFunction(fns)) {
    router.add(category, fns);
    return;
  }
  _.forEach(fns, (v, k) => {
    if (_.isFunction(v)) {
      debug('init route:%s', `${category}.${k}`);
      router.add(`${category}.${k}`, v);
    } else if (_.isObject(v)) {
      addToRouter(`${category}.${k}`, v);
    } else {
      /* istanbul ignore next */
      console.error(`${category}.${k} is invalid.`);
    }
  });
}

router.addDefault('common', middlewares.common.routeStats());

addToRouter('c', controllers);
addToRouter('m.noQuery', middlewares.common.noQuery());
addToRouter('m.noCache', middlewares.common.noCache());
addToRouter('m.session', middlewares.session.writable());
addToRouter('m.anonymous', middlewares.session.anonymous());
addToRouter('m.login', middlewares.session.login());
addToRouter('m.admin', middlewares.session.admin());
addToRouter('m.su', middlewares.session.su());

addToRouter('m.level', middlewares.level);
addToRouter('m.version', middlewares.common.version);
addToRouter('m.tracker', middlewares.tracker);

export default getRouter(routes);
