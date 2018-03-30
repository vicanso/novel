import Joi from 'joi';
import _ from 'lodash';

import influx from '../helpers/influx';
import errors from '../errors';

const schema = {
  measurement: () => Joi.string().valid(['userTracker']),
  limit: () =>
    Joi.number()
      .integer()
      .min(1)
      .max(1000)
      .default(100),
  start: () =>
    Joi.string()
      .max(32)
      .default('-1d'),
  end: () => Joi.string().max(32),
  conditions: () => Joi.string().max(200),
};

// 获取series
export async function getSeries(ctx) {
  const measurement = Joi.attempt(ctx.params.measurement, schema.measurement());
  const client = influx.getClient();
  if (!client) {
    throw errors.get('common.influxNotInit');
  }
  const data = await client.showSeries(measurement);
  const tags = {};
  _.forEach(data, item => {
    const arr = item.split(',');
    arr.shift();
    _.forEach(arr, subItem => {
      const [key, value] = subItem.split('=');
      if (!tags[key]) {
        tags[key] = [];
      }
      tags[key].push(value);
    });
  });
  _.forEach(tags, (value, key) => {
    tags[key] = _.uniq(value).sort();
  });
  ctx.setCache('5m');
  ctx.body = tags;
}

// 查询measurement符合条件的数据
export async function list(ctx) {
  const measurement = Joi.attempt(ctx.params.measurement, schema.measurement());
  const {limit, start, end, conditions} = Joi.validate(ctx.query, {
    limit: schema.limit(),
    start: schema.start(),
    end: schema.end(),
    conditions: schema.conditions(),
  });
  const client = influx.getClient();
  if (!client) {
    throw errors.get('common.influxNotInit');
  }
  const ql = client.query(measurement);
  if (start) {
    ql.start = start;
  }
  if (end) {
    ql.end = end;
  }
  if (conditions) {
    const conds = JSON.parse(conditions);
    _.forEach(conds, (v, k) => ql.condition(k, v));
  }
  const data = await ql.set({
    limit,
    format: 'json',
  });
  ctx.body = {
    trackers: data[measurement],
  };
}
