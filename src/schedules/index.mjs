import schedule from 'node-schedule';

import performance from './performance';
import {isDevelopment} from '../helpers/utils';
import {updateAll} from '../services/book';

performance(10 * 1000);

schedule.scheduleJob('00 * * * *', () => {
  if (isDevelopment()) {
    return;
  }
  updateAll();
});
