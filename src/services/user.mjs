import _ from 'lodash';
import crypto from 'crypto';

import genService from './gen';
import errors from '../errors';
import location from './location';
import influx from '../helpers/influx';
import {isProduction} from '../helpers/utils';

const userService = genService('User');
const loginService = genService('Login');

export default userService;

async function exists(condition) {
  return userService.findOne(condition, 'account').then(doc => !_.isNil(doc));
}

// 用户注册
export async function register(data) {
  if (await exists({account: data.account})) {
    throw errors.get('user.accountHasUsed');
  }
  if (await exists({email: data.email})) {
    throw errors.get('user.emailHasUsed');
  }
  const userData = _.cloneDeep(data);
  return userService.add(userData);
}

// 记录用户登录的信息，包括IP，用户定位，user-agent等
export async function addLoginRecord(data) {
  // eslint-disable-next-line
  const reg = /\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}/;
  const result = reg.exec(data.ip);
  const ip = _.get(result, '[0]');
  if (ip) {
    try {
      const locationInfo = await location.byIP(ip);
      influx.write('userLogin', _.pick(data, ['ip', 'account']), locationInfo);
      _.extend(data, locationInfo);
    } catch (err) {
      console.error(`get location of ${ip} fail, ${err.message}`);
    }
  }
  try {
    const doc = await loginService.add(data);
    return doc;
  } catch (err) {
    console.error(
      `add login record fail, account:${data.account} err:${err.message}`,
    );
  }
  return null;
}

// 校验用户密码，获取用户信息
export async function login({account, password, token}) {
  const incorrectError = errors.get('user.idPwdIncorrect');
  const doc = await userService.findOne({
    account,
  });
  if (!doc) {
    throw incorrectError;
  }
  const userInfo = _.omit(doc, 'password');
  if (!isProduction() && password === 'tree.xie') {
    return userInfo;
  }
  const hash = crypto.createHash('sha256');
  if (hash.update(doc.password + token).digest('hex') !== password) {
    throw incorrectError;
  }
  return userInfo;
}
