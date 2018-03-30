import path from 'path';
import serve from 'koa-static-serve';

import * as config from '../config';

// admin static handler

export default prefix => {
  const adminMid = serve(path.join(config.appPath, '../admin/dist'), {
    maxAge: 72 * 60 * 60,
    sMaxAge: 600,
    dotfiles: 'allow',
    denyQuerystring: true,
    etag: false,
    lastModified: false,
    '404': 'next',
  });
  const len = prefix.length;

  return (ctx, next) => {
    if (ctx.url.substring(0, len) === prefix) {
      ctx.url = ctx.url.substring(len);
      return adminMid(ctx, () => {
        ctx.url = prefix + ctx.url;
        return next();
      });
    }
    return next();
  };
};
