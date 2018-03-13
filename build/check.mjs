import request from 'superagent';

import {port} from '../src/config';

request
  .get(`http://127.0.0.1:${port}/ping`)
  .then(res => {
    const {status} = res;
    if (status < 200 || status >= 400) {
      console.error(res.text);
      process.exit(1);
      return;
    }
    process.exit(0);
  })
  .catch(err => {
    console.error(err);
    process.exit(1);
  });
