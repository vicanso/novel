import Joi from 'joi';
import _ from 'lodash';

import mockService from '../services/mock';

/**
 * @swagger
 * parameters:
 *  mock.account:
 *    name: account
 *    description: 需要Mock的账号，Joi.string().trim().max(100)
 *    type: string
 *    in: formData
 *  mock.url:
 *    name: url
 *    description: 需要Mock的url，Joi.string().trim().max(200)
 *    type: string
 *    in: formData
 *  mock.status:
 *    name: status
 *    description: 响应的状态码
 *    type: integer
 *    in: formData
 *  mock.response:
 *    name: response
 *    description: 响应的数据
 *    type: object
 *    in: formData
 *  mock.description:
 *    name: description
 *    description: 该配置描述
 *    type: string
 *    in: formData
 *  mock.disabled:
 *    name: disabled
 *    description: 是否禁用
 *    type: boolean
 *    in: formData
 */
const schema = {
  account: () =>
    Joi.string()
      .trim()
      .max(100),
  url: () =>
    Joi.string()
      .trim()
      .max(200),
  status: () => Joi.number().integer(),
  response: () => Joi.object(),
  description: () =>
    Joi.string()
      .trim()
      .max(300),
  disabled: () => Joi.boolean(),
};

/**
 * @swagger
 * definitions:
 *  Mock:
 *    required:
 *      - url
 *    properties:
 *      account:
 *        description: 账号
 *        type: string
 *      url:
 *        description: url
 *        type: string
 *      status:
 *        description: 响应状态码
 *        type: integer
 *      response:
 *        description: 响应数据
 *        type: object
 *      description:
 *        description: 该mock描述
 *        type: string
 *      createdAt:
 *        description: 创建时间
 *        type: string
 *      creator:
 *        description: 创建者
 *        type: string
 *      disabled:
 *        description: 是否禁用
 *        type: boolean
 */

/**
 * @swagger
 * /mocks:
 *  post:
 *    description: 创建新的Mock配置。中间件：m.admin
 *    summary: 创建新的Mock配置
 *    tags:
 *      - mock
 *    consumes:
 *      - multipart/form-data
 *    parammeters:
 *      - $ref: '#/parameters/mock.account'
 *      - $ref: '#/parameters/mock.url'
 *      - $ref: '#/parameters/mock.status'
 *      - $ref: '#/parameters/mock.response'
 *      - $ref: '#/parameters/mock.description'
 *      - $ref: '#/parameters/mock.disabled'
 *    responses:
 *      201:
 *        description: 创建成功
 *        schema:
 *          type: object
 *          $ref: '#/definitions/Mock'
 */
export async function add(ctx) {
  const data = Joi.validate(ctx.request.body, {
    account: schema.account(),
    url: schema.url().required(),
    status: schema.status(),
    response: schema.response(),
    description: schema.description(),
    disabled: schema.disabled(),
  });
  data.creator = ctx.session.user.account;
  const doc = await mockService.add(data);
  ctx.status = 201;
  ctx.body = doc;
}

/**
 * @swagger
 * /mocks:
 *  get:
 *    description: 获取Mock配置列表。中间件：m.noQuery
 *    summary: 获取Mock配置列表
 *    tags:
 *      - mock
 *    responses:
 *      200:
 *        description: 返回Mock配置列表
 *        schema:
 *          type: object
 *          properties:
 *            list:
 *              type: array
 *              items:
 *                $ref: '#/definitions/Mock'
 */
export async function list(ctx) {
  const mocks = await mockService.find({});
  ctx.setCache('5s');
  ctx.body = {
    list: _.sortBy(mocks, item => item.url),
  };
}

/**
 * @swagger
 * /mocks/:id:
 *  get:
 *    description: 获取该Mock配置。中间件：m.noQuery
 *    summary: 获取该Mock配置
 *    tags:
 *      - mock
 *    parammeters:
 *      - $ref: '#/parameters/mock.account'
 *      - $ref: '#/parameters/mock.url'
 *      - $ref: '#/parameters/mock.status'
 *      - $ref: '#/parameters/mock.response'
 *      - $ref: '#/parameters/mock.description'
 *      - $ref: '#/parameters/mock.disabled'
 *    responses:
 *      200:
 *        description: 返回该Mock配置
 *        schema:
 *          type: object
 *          $ref: '#/definitions/Mock'
 */
export async function get(ctx) {
  const id = Joi.attempt(ctx.params.id, Joi.objectId());
  const doc = await mockService.findById(id);
  ctx.setCache('5s');
  ctx.body = doc;
}

/**
 * @swagger
 * /mocks/:id:
 *  patch:
 *    description: 更新Mock配置。中间件：m.admin
 *    summary: 更新Mock配置
 *    tags:
 *      - mock
 *    consumes:
 *      - multipart/form-data
 *    responses:
 *      204:
 *        description: 更新成功
 */
export async function update(ctx) {
  const data = Joi.validate(ctx.request.body, {
    account: schema.account(),
    url: schema.url(),
    status: schema.status(),
    response: schema.response(),
    description: schema.description(),
    disabled: schema.disabled(),
  });
  const id = Joi.attempt(ctx.params.id, Joi.objectId());
  await mockService.findByIdAndUpdate(id, data);
  ctx.body = null;
}
