import Vue from "vue";
import Router from "vue-router";

import { routeDetail, routeLogin, routeRegister } from "@/routes";
import Detail from "@/views/Detail";
import Login from "@/views/Login";
import Register from "@/views/Register";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/detail/:id",
      name: routeDetail,
      component: Detail
    },
    {
      path: "/login",
      name: routeLogin,
      component: Login
    },
    {
      path: "/register",
      name: routeRegister,
      component: Register
    }
  ]
});
