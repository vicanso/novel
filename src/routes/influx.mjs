export default [
  ['GET', '/influx/:measurement', 'm.admin & c.influx.list'],
  ['GET', '/influx/series/:measurement', 'm.noQuery & c.influx.getSeries'],
];
