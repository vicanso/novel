import schedule from 'node-schedule';

import performance from './performance';
import {isDevelopment} from '../helpers/utils';
import {updateAll, updateAllCategory} from '../services/book';

performance(10 * 1000);
schedule.scheduleJob('00 * * * *', async () => {
  if (isDevelopment()) {
    return;
  }
  try {
    await updateAll();
    console.info('schedule update all completed');
  } catch (err) {
    console.error(`schedule update all fail, ${err.message}`);
  }
});

schedule.scheduleJob('30 00 * * *', async () => {
  if (isDevelopment()) {
    return;
  }
  try {
    await updateAllCategory();
    console.info('schedule update all category completed');
  } catch (err) {
    console.error(`schedule update all category fail, ${err.message}`);
  }
});
