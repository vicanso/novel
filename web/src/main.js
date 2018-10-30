import bluebird from "bluebird";
global.Promise = bluebird;
import "intersection-observer";
import Vue from "vue";
import App from "@/App.vue";
import router from "@/router";
import store from "@/store";
import Mint from "mint-ui";
import "@/styles/index.scss";
import "@/request-interceptors";
import { getErrorMessage, isDevelopment } from "@/helpers/util";
import { timeout } from "@/config";
import { clearChapterStoreExpired } from "@/helpers/storage";

import "@/assets/iconfont/iconfont.css";
import "@/directive";

Vue.use(Mint);

// 注入 router 和 store
Vue.$router = router;
Vue.$store = store;

Vue.prototype.xLoading = (options = {}) => {
  const Indicator = Mint.Indicator;
  Indicator.open({
    text: "加载中...",
    spinnerType: "fading-circle"
  });
  let resolved = false;
  const resolve = () => {
    if (resolved) {
      return;
    }
    resolved = true;
    Indicator.close();
  };
  setTimeout(resolve, options.timeout || timeout);
  return resolve;
};

Vue.prototype.xError = function xError(err) {
  const message = getErrorMessage(err);
  Mint.Toast({
    message: message,
    duration: 5000
  });
  if (isDevelopment()) {
    throw err;
  }
};
Vue.prototype.xToast = (message, ms = 3000) => {
  Mint.Toast({
    message,
    duration: ms
  });
};

Vue.prototype.$next = function nextTickPromise() {
  return new Promise(resolve => this.$nextTick(resolve));
};

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");

setTimeout(() => {
  clearChapterStoreExpired();
}, 5000);
