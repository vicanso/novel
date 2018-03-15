import schedule from 'node-schedule';

import performance from './performance';
import {isDevelopment} from '../helpers/utils';
import {updateAll} from '../services/book';

performance(10 * 1000);
schedule.scheduleJob('10 * * * *', () => {
  if (isDevelopment()) {
    return;
  }
  updateAll()
    .then(() => {
      console.info('schedule update all completed');
    })
    .catch(err => {
      console.error(`schedule update all fail, ${err.message}`);
    });
});
