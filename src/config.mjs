import shortid from 'shortid';

import pkg from '../package';
import expose from './expose';

export const port = Number.parseInt(process.env.PORT, 10) || 5018;

export const env = process.env.NODE_ENV || 'development';

// eslint-disable-next-line
export const version = pkg.version;

export const app = pkg.name;

// eslint-disable-next-line
export const httpLogFormat =
  ':account :request-id :method :url :status :length :response-time ms ":referrer"';

export const trackCookie = 'jt';

export const appUrlPrefix = '/api';

export const ins = shortid();

// logger setting "console", "udp://127.0.0.1:5001"
export const logger = process.env.LOG;

// mongodb connection uri
export const mongoUri =
  process.env.MONGO || 'mongodb://127.0.0.1/novel?connectTimeoutMS=300000';

// redis connection uri
export const redisUri = process.env.REDIS || 'redis://127.0.0.1/';

export const appPath = expose.dirname;

// cookie name
export const session = {
  key: pkg.name,
  maxAge: 1 * 24 * 3600 * 1000,
};

// http connection limit options
export const connectLimitOptions = {
  mid: 100,
  high: 500,
  interval: 5000,
};
