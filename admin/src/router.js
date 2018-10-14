import Vue from "vue";
import Router from "vue-router";
import Home from "@/views/Home.vue";
import Login from "@/views/Login.vue";
import Register from "@/views/Register.vue";
import BookList from "@/views/BookList.vue";
import { routeLogin, routeHome, routeRegister, routeBookList } from "@/routes";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: routeHome,
      component: Home
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
    },
    {
      path: "/book-list",
      name: routeBookList,
      component: BookList
    }
  ]
});
