import Vue from "vue";
import Vuex from "vuex";

import userStore from "@/store/user";
import bookStore from "@/store/book";

Vue.use(Vuex);

const defaultOptions = {
  state: {},
  mutations: {},
  actions: {}
};
[userStore, bookStore].forEach(item => {
  Object.assign(defaultOptions.state, item.state);
  Object.assign(defaultOptions.mutations, item.mutations);
  Object.assign(defaultOptions.actions, item.actions);
});

export default new Vuex.Store(defaultOptions);
