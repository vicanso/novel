import Joi from 'joi';
import _ from 'lodash';
import util from 'util';
import fs from 'fs';
import path from 'path';

import i18nService from '../services/i18n';
import debug from '../helpers/debug';
import * as config from '../config';

const readFile = util.promisify(fs.readFile);

/**
 * @swagger
 * parameters:
 *  i18n.name:
 *    name: name
 *    description: 配置名称，Joi.string().trim().max(100)
 *    type: string
 *    in: formData
 *  i18n.category:
 *    name: category
 *    description: 配置类型，Joi.string().trim().max(20)
 *    type: string
 *    in: formData
 *  i18n.en:
 *    name: en
 *    description: 英文，Joi.string().trim().max(200)
 *    type: string
 *    in: formData
 *  i18n.zh:
 *    name: zh
 *    description: 中文，Joi.string().trim().max(100),
 *    type: string
 *    in: formData
 */
const getDefaultSchema = () => ({
  name: Joi.string()
    .trim()
    .max(100),
  category: Joi.string()
    .trim()
    .max(20),
  en: Joi.string()
    .trim()
    .max(200),
  zh: Joi.string()
    .trim()
    .max(100),
});

/**
 * @swagger
 * definitions:
 *  I18n:
 *    required:
 *      - name
 *      - category
 *      - createdAt
 *      - creator
 *    properties:
 *      name:
 *        description: 配置名称
 *        type: string
 *        example: login
 *      category:
 *        description: 配置分类
 *        type: string
 *        example: basic
 *      en:
 *        description: 英文
 *        type: string
 *        example: login
 *      zh:
 *        description: 中文
 *        type: string
 *        example: 登录
 *      createdAt:
 *        description: 创建时间，ISO Date String
 *        type: string
 *      creator:
 *        description: 创建者
 *        type: string
 *        example: vicanso
 */

/**
 * @swagger
 * /i18ns:
 *  post:
 *    description: 创建新的多语言配置。中间件：m.admin
 *    summary: 创建新的多语言配置
 *    tags:
 *      - i18n
 *    consumes:
 *      - multipart/form-data
 *    parameters:
 *      - $ref: '#/parameters/i18n.name'
 *      - $ref: '#/parameters/i18n.category'
 *      - $ref: '#/parameters/i18n.en'
 *      - $ref: '#/parameters/i18n.zh'
 *    responses:
 *      201:
 *        description: 创建成功
 *        schema:
 *          type: object
 *          $ref: '#/definitions/I18n'
 */
export async function add(ctx) {
  const schema = getDefaultSchema();
  schema.name.required();
  schema.category.required();
  const data = Joi.validate(ctx.request.body, schema);
  data.creator = ctx.session.user.account;
  const doc = await i18nService.add(data);
  ctx.status = 201;
  ctx.body = doc;
}

/**
 * @swagger
 * /i18ns/:id:
 *  patch:
 *    description: 更新多语言配置。中间件：m.admin
 *    summary: 更新多语言配置
 *    tags:
 *      - i18n
 *    consumes:
 *      - multipart/form-data
 *    parameters:
 *      - $ref: '#/parameters/i18n.name'
 *      - $ref: '#/parameters/i18n.category'
 *      - $ref: '#/parameters/i18n.en'
 *      - $ref: '#/parameters/i18n.zh'
 *    responses:
 *      204:
 *        description: 更新成功
 */
export async function update(ctx) {
  const data = Joi.validate(ctx.request.body, getDefaultSchema());
  const id = Joi.attempt(ctx.params.id, Joi.objectId());
  await i18nService.findByIdAndUpdate(id, data);
  ctx.body = null;
}

/**
 * @swagger
 * /i18ns:
 *  get:
 *    description: 获取多语言配置
 *    summary: 获取多语言配置
 *    tags:
 *      - i18n
 *    parameters:
 *      -
 *        name: category
 *        description: 需要获取的分类，可以一次获取多个分类
 *        type: string
 *        in: query
 *      -
 *        name: count
 *        description: 是否获取总数
 *        type: boolean
 *        in: query
 *    responses:
 *      200:
 *        description: 返回配置列表
 *        schema:
 *          type: object
 *          properties:
 *            list:
 *              type: array
 *              items:
 *                $ref: '#/definitions/I18n'
 *            count:
 *              type: integer
 *              example: 100
 */
export async function list(ctx) {
  const schema = getDefaultSchema();
  const {category, count} = Joi.validate(ctx.query, {
    category: Joi.alternatives().try(
      schema.category,
      Joi.array().items(schema.category),
    ),
    count: Joi.boolean(),
  });
  let cats = category;
  const conditions = {};
  if (cats) {
    if (!_.isArray(cats)) {
      cats = [cats];
    }
    conditions.category = cats;
  }
  const items = await i18nService.find(conditions);
  const result = {
    list: _.sortBy(items, item => item.category + item.name),
  };
  if (count) {
    result.count = await i18nService.count(conditions);
  }
  ctx.setCache('5s');
  debug('list i18n conditions:%j, result:%j', conditions, result);
  ctx.body = result;
}

/**
 * @swagger
 * /i18ns/categories:
 *  get:
 *    description: 获取i18n所有分类。中间件：m.noQuery
 *    summary: 获取i18n分类
 *    tags:
 *      - i18n
 *    responses:
 *      200:
 *        description: 返回分类列表
 *        schema:
 *          type: object
 *          properties:
 *            list:
 *              type: array
 *              example: ["basic", "user"]
 *              items:
 *                type: string
 */
export async function listCategory(ctx) {
  const items = await i18nService.find({}).select('category');
  const categories = [];
  _.forEach(items, item => {
    if (!_.includes(categories, item.category)) {
      categories.push(item.category);
    }
  });
  ctx.setCache('5m');
  ctx.body = {
    list: categories,
  };
}

/**
 * @swagger
 * /i18ns/category:
 *  get:
 *    description: 根据选择语言，获取分类的配置信息
 *    summary: 获取分类的配置信息
 *    tags:
 *      - i18n
 *    parameters:
 *      -
 *        name: category
 *        description: 需要获取的分类，可以一次获取多个
 *        type: string
 *        in: query
 *        required: true
 *      -
 *        name: lang
 *        description: 语言
 *        type: string
 *        in: query
 *    responses:
 *      200:
 *        description: 成功返回该类型该语言的配置
 */
export async function listCategoryByLang(ctx) {
  const schema = getDefaultSchema();
  const {category} = Joi.validate(ctx.query, {
    category: Joi.alternatives()
      .try(schema.category, Joi.array().items(schema.category))
      .required(),
  });
  let cats = category;
  const {lang} = ctx.state;
  const conditions = {};
  if (!_.isArray(cats)) {
    cats = [cats];
  }
  conditions.category = cats;
  const items = await i18nService.find(conditions);
  const result = {};
  _.forEach(items, item => {
    const cat = item.category;
    if (!result[cat]) {
      result[cat] = {};
    }
    result[cat][item.name] = item[lang];
  });
  debug('select langs of %s, result:%j', category, result);
  ctx.setCache('10m', '1m');
  ctx.body = result;
}

/**
 * @swagger
 * /i18ns/:id:
 *  get:
 *    description: 获取该多语言配置。中间件：m.noQuery
 *    summary: 获取该多语言配置
 *    tags:
 *      - i18n
 *    responses:
 *      200:
 *        description: 返回该多语言配置
 *        schema:
 *          type: object
 *          $ref: '#/definitions/I18n'
 */
export async function get(ctx) {
  const id = Joi.attempt(ctx.params.id, Joi.objectId());
  const doc = await i18nService.findById(id);
  ctx.setCache('10s');
  ctx.body = doc;
}

// 初始化i18n配置信息
export async function init(ctx) {
  if (ctx.request.body.token !== 'SJhZZwZ0b') {
    throw new Error('Token is invalid');
  }
  const buf = await readFile(path.join(config.appPath, 'assets/i18n.json'));
  _.forEach(JSON.parse(buf), async item => {
    // eslint-disable-next-line
    item.creator = 'vicanso';
    const conditions = _.pick(item, ['category', 'name']);
    await i18nService.findOneAndUpdate(conditions, item, {
      upsert: true,
    });
  });
  ctx.status = 201;
}
