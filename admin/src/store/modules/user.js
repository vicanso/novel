import request from 'axios';
import _ from 'lodash';

import {sha256} from '../../helpers/crypto';
import {
  USERS_ME,
  USERS_LOGIN,
} from '../../urls';
import {
  USER_INFO,
} from '../mutation-types';
import {genPassword} from '../../helpers/util';


const state = {
  info: null,
};

const mutations = {
  // 用户信息
  [USER_INFO](state, data) {
    const isAdmin = _.includes(data.roles, 'admin');
    state.info = _.extend(
      {
        isAdmin,
      },
      data,
    );
  },
};

const userGetInfo = async ({commit}) => {
  const res = await request.get(USERS_ME, {
    params: {
      'cache-control': 'no-cache',
    },
  });
  commit(USER_INFO, res.data);
};

// 用户注册
const userRegister = async ({commit}, {account, password, email}) => {
  const res = await request.post(USERS_ME, {
    account,
    password: genPassword(account, password),
    email,
  });
  commit(USER_INFO, res.data);
};

// 用户登录
const userLogin = async ({commit}, {account, password}) => {
  let res = await request.get(USERS_LOGIN, {
    params: {
      'cache-control': 'no-cache',
    },
  });
  const token = res.data.token;
  const code = sha256(genPassword(account, password) + token);
  res = await request.post(USERS_LOGIN, {
    account,
    password: code,
  });
  commit(USER_INFO, res.data);
};

const userLogout = async ({commit}) => {
  await request.delete(USERS_ME);
  commit(USER_INFO, {
    anonymous: true,
  });
};


export const actions = {
  userGetInfo,
  userRegister,
  userLogin,
  userLogout,
};

export default {
  state,
  mutations,
};