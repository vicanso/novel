import * as userModule from './user';
import * as bookModule from './book';

const modules = {};
const actions = {};
const getters = {};

modules.book = bookModule.default;
Object.assign(actions, bookModule.actions);

modules.user = userModule.default;
Object.assign(actions, userModule.actions);

export default {
  actions,
  modules,
  getters,
};