import Logger from 'timtam-logger';
import als from 'async-local-storage';

import * as config from '../config';

if (config.logger) {
  const logger = new Logger({
    app: config.app,
  });
  logger.before(() => als.get('account'));
  logger.before(() => als.get('id'));
  logger.wrap(console);
  'emerg alert crit'.split(' ').forEach(event => {
    logger.on(event, message => {
      console.dir(message);
      // TODO 发送email警报
    });
  });
  logger.add(config.logger);
} else {
  console.emerg = console.error;
  console.alert = console.error;
  console.crit = console.error;
}

export default console;
