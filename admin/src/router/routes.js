import Setting from '../views/setting';
import Register from '../views/register';
import Login from '../views/login';
import Mock from '../views/mock';
import Tracker from '../views/tracker';
import Book from '../views/book';

export default [
  {
    name: 'register',
    path: '/register',
    component: Register,
  },
  {
    name: 'setting',
    path: '/setting',
    component: Setting,
  },
  {
    name: 'login',
    path: '/login',
    component: Login,
  },
  {
    name: 'mock',
    path: '/mock',
    component: Mock,
  },
  {
    name: 'tracker',
    path: '/tracker',
    component: Tracker,
  },
  {
    name: 'book',
    path: '/book',
    component: Book,
  },
];
