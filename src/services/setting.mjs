import _ from 'lodash';

import genService from './gen';

const settingService = genService('Setting');

const {find} = settingService;

let applicationSettings = null;

export default settingService;

export function get(key) {
  return _.get(applicationSettings, key);
}

/**
 * 更新应用配置信息
 *
 */
export async function updateAppSettings() {
  const docs = await find({
    disabled: false,
  });
  applicationSettings = {};
  _.forEach(docs, item => {
    applicationSettings[item.name] = item.data;
  });
}
