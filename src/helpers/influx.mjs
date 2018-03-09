import _ from 'lodash';
import Influx from 'influxdb-nodejs';
import stringify from 'simple-stringify';

import schemas from '../influx-schemas';
import debug from './debug';
import {getParam} from './utils';
import * as config from '../config';

let client = null;
const maxQueueLength = 100;

/**
 * 将当前队列中的统计数据全部写入到influxdb
 */
function flush() {
  const count = client.writeQueueLength;
  if (!count) {
    return;
  }
  client
    .syncWrite()
    .then(() => console.info(`influxdb write ${count} records sucess`))
    .catch(err => console.error(`influxdb write fail, ${err.message}`));
}

function init(url) {
  client = new Influx(url);
  const debounceFlush = _.debounce(flush, 30 * 1000);
  Object.keys(schemas).forEach(measurement => {
    const {fields, tags, options} = schemas[measurement];
    client.schema(measurement, fields, tags, options);
  });
  client.timeout = 3000;
  client.on('writeQueue', () => {
    // sync write queue if the length is 100
    if (client.writeQueueLength >= maxQueueLength) {
      flush();
    } else {
      debounceFlush();
    }
  });
  client.on('queue', ({type, data}) => {
    // 全局增加server field
    if (type === 'write') {
      // eslint-disable-next-line
      data.fields.ins = config.ins;
    }
  });

  client.on('invalid-fields', data => {
    console.error(`influx invalid fields:${stringify.json(data, 3)}`);
  });
  client.on('invalid-tags', data => {
    console.error(`influx invalid tags:${stringify.json(data, 3)}`);
  });
}

function write(measurement, fields, ...args) {
  /* istanbul ignore if */
  if (!client) {
    debug('measurement:%s, fields:%j, args:%j', measurement, fields, args);
    return null;
  }
  const writer = client.write(measurement).field(fields);
  const tags = getParam(args, _.isObject);
  /* istanbul ignore else */
  if (tags) {
    writer.tag(tags);
  }
  debug('influx measurement:%s, fields:%j, tags:%j', measurement, fields, tags);
  const syncNow = getParam(args, _.isBoolean);
  if (!syncNow) {
    writer.queue();
  }
  return writer;
}

export default {
  init,
  getClient: () => client,
  write,
};
