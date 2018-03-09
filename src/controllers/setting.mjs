import Joi from 'joi';
import _ from 'lodash';

import settingService from '../services/setting';

/**
 * @swagger
 *  parameters:
 *    setting.name:
 *      name: name
 *      description: 该配置名称
 *      type: string
 *      in: formData
 *    setting.data:
 *      name: data
 *      description: 该配置数据
 *      type: object
 *      in: formData
 *    setting.disabled:
 *      name: disabled
 *      description: 该配置是否禁用
 *      type: boolean
 *      in: formData
 *    setting.description:
 *      name: description
 *      description: 该配置的描述
 *      type: string
 *      in: formData
 */
const getDefaultSchema = () => ({
  name: Joi.string()
    .trim()
    .max(30),
  data: Joi.object(),
  disabled: Joi.boolean(),
  description: Joi.string().trim(),
});

/**
 * @swagger
 * definitions:
 *  Setting:
 *    required:
 *      - name
 *    properties:
 *      name:
 *        description: 配置名称
 *        type: string
 *      createdAt:
 *        description: 创建时间
 *        type: string
 *      disabled:
 *        description: 配置是否禁用
 *        type: boolean
 *      data:
 *        description: 配置数据
 *        type: object
 *      creator:
 *        description: 创建者
 *        type: string
 *      description:
 *        description: 配置描述
 *        type: string
 */

/**
 * @swagger
 * /settings:
 *  post:
 *    description: 创建新的Setting配置。中间件：m.admin
 *    summary: 创建新的Setting配置
 *    tags:
 *      - setting
 *    consumes:
 *      - multipart/form-data
 *    parammeters:
 *      - $ref: '#/parameters/setting.name'
 *      - $ref: '#/parameters/setting.data'
 *      - $ref: '#/parameters/setting.disabled'
 *      - $ref: '#/parameters/setting.description'
 *    responses:
 *      201:
 *        description: 创建成功
 *        schema:
 *          type: object
 *          $ref: '#/definitions/Setting'
 */
export async function add(ctx) {
  const schema = getDefaultSchema();
  schema.name.required();
  const data = Joi.validate(ctx.request.body, schema);
  data.creator = ctx.session.user.account;
  const doc = await settingService.add(data);
  ctx.status = 201;
  ctx.body = doc;
}

/**
 * @swagger
 * /settings:
 *  get:
 *    description: 获取所有的配置信息。中间件：m.admin
 *    summary: 获取配置信息
 *    tags:
 *      - setting
 *    responses:
 *      200:
 *        description: 返回Setting信息列表
 *        schema:
 *          type: object
 *          properties:
 *            list:
 *              type: array
 *              items:
 *                $ref: '#/definitions/Setting'
 */
export async function list(ctx) {
  const items = await settingService.find({});
  ctx.body = {
    list: items,
  };
}

/**
 * @swagger
 * /settings/:id:
 *  get:
 *    description: 获取该Setting信息。中间件：m.admin
 *    summary: 获取该Setting信息
 *    tags:
 *      - setting
 *    responses:
 *      200:
 *        description: 返回该Setting信息
 *        $ref: '#/definitions/Setting'
 */
export async function get(ctx) {
  const id = Joi.attempt(ctx.params.id, Joi.objectId());
  const doc = await settingService.findById(id);
  ctx.body = doc;
}

/**
 * @swagger
 * /settings/:id:
 *  patch:
 *    description: 更新setting配置。中间件：m.admin
 *    summary: 更新setting配置
 *    tags:
 *      - setting
 *    consumes:
 *      - multipart/form-data
 *    parammeters:
 *      - $ref: '#/parameters/setting.name'
 *      - $ref: '#/parameters/setting.data'
 *      - $ref: '#/parameters/setting.disabled'
 *      - $ref: '#/parameters/setting.description'
 *    responses:
 *      204:
 *        description: 更新成功
 */
export async function update(ctx) {
  const id = Joi.attempt(ctx.params.id, Joi.objectId());
  const data = Joi.validate(ctx.request.body, getDefaultSchema());
  if (!_.isEmpty(data)) {
    await settingService.findByIdAndUpdate(id, data);
  }
  ctx.body = null;
}
