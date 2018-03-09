import _ from 'lodash';

import sys from './system';
import i18n from './i18n';
import user from './user';
import mock from './mock';
import setting from './setting';
import stats from './stats';

export default _.flatten([sys, i18n, user, mock, setting, stats]);
