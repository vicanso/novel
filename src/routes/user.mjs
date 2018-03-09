export default [
  ['GET', '/users/me', 'm.session & c.user.me'],
  ['PATCH', '/users/me', 'm.session & c.user.refresh'],
  [
    'POST',
    '/users/register',
    ['m.anonymous', 'm.tracker("register")', 'c.user.register'],
  ],
  [
    'DELETE',
    '/users/logout',
    ['m.login', 'm.tracker("logout")', 'c.user.logout'],
  ],
  ['GET', '/users/login', 'm.anonymous & c.user.loginToken'],
  [
    'POST',
    '/users/login',
    ['m.anonymous', 'm.tracker("login")', 'c.user.login'],
  ],
];
