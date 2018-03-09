import _ from 'lodash';
import stringify from 'simple-stringify';

/**
 * @swagger
 * /stats/requests:
 *  post:
 *    description: 提交接口请求相关系统
 *    responses:
 *      201:
 *        description: 记录成功
 */
export function add(ctx) {
  _.forEach(ctx.request.body, item => {
    console.info(`browser-ajax ${stringify.json(item)}`);
  });
  ctx.status = 201;
}

/**
 * @swagger
 * /stats/exceptions:
 *  post:
 *    description: 提交出错信息统计
 *    responses:
 *      201:
 *        description: 记录成功
 */
export function exception(ctx) {
  _.forEach(ctx.request.body, item => {
    console.error(`browser-exception ${stringify.json(item)}`);
  });
  ctx.status = 201;
}

/**
 * @swagger
 * /stats/routes:
 *  post:
 *    description: 提交route统计信息
 *    responses:
 *      201:
 *        description: 记录成功
 */
export function route(ctx) {
  _.forEach(ctx.request.body, item => {
    console.info(`route ${stringify.json(item)}`);
  });
  ctx.status = 201;
}
