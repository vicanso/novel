import Joi from 'joi';
import Promise from 'bluebird';
import zlib from 'zlib';
import util from 'util';

import orginService from '../services/origin';
import bookService, {addBook} from '../services/book';
import chapterService from '../services/chapter';
import coverService from '../services/cover';

const gunzip = util.promisify(zlib.gunzip);

// 增加来源
export async function addSource(ctx) {
  const {name, author, source, id} = Joi.validate(ctx.request.body, {
    name: Joi.string()
      .max(30)
      .required(),
    author: Joi.string()
      .max(20)
      .required(),
    source: Joi.string()
      .valid(['biquge'])
      .required(),
    id: Joi.string()
      .max(10)
      .required(),
  });
  const doc = await orginService.add({
    name,
    author,
    source,
    sourceId: id,
  });
  ctx.status = 201;
  ctx.body = doc;
}

// 获取书籍列表
export async function list(ctx) {
  const {skip, limit, count, fields} = Joi.validate(ctx.query, {
    skip: Joi.number()
      .integer()
      .default(0),
    limit: Joi.number()
      .integer()
      .min(1)
      .max(20)
      .default(10),
    fields: Joi.string().max(100),
    count: Joi.boolean(),
  });
  const conditions = {};
  const data = {};
  if (count) {
    data.count = await bookService.count(conditions);
  }
  data.list = await bookService
    .find(conditions)
    .skip(skip)
    .limit(limit)
    .select(fields)
    .lean();
  ctx.setCache(60);
  ctx.body = data;
}

// 增加书籍
export async function add(ctx) {
  const {name, author} = Joi.validate(ctx.request.body, {
    name: Joi.string()
      .max(30)
      .required(),
    author: Joi.string()
      .max(20)
      .required(),
  });
  const doc = await addBook(author, name);
  ctx.status = 201;
  ctx.body = doc;
}

// 获取书籍信息
export async function get(ctx) {
  const no = Joi.attempt(
    ctx.params.no,
    Joi.number()
      .integer()
      .min(0)
      .max(10000)
      .required(),
  );
  const doc = await bookService.findOne({
    no,
  });
  ctx.body = doc;
}

// 列出章节信息
export async function listChapter(ctx) {
  const no = Joi.attempt(
    ctx.params.no,
    Joi.number()
      .integer()
      .min(0)
      .max(10000)
      .required(),
  );
  const {skip, limit, fields, sort} = Joi.validate(ctx.query, {
    skip: Joi.number()
      .integer()
      .default(0),
    limit: Joi.number()
      .integer()
      .min(1)
      .max(100)
      .default(10),
    fields: Joi.string().max(100),
    sort: Joi.string()
      .max(30)
      .default('updatedAt'),
  });
  const docs = await chapterService
    .find({
      book: no,
    })
    .sort(sort)
    .skip(skip)
    .limit(limit)
    .select(fields)
    .lean();
  await Promise.mapSeries(docs, async item => {
    if (item.data) {
      const content = await gunzip(item.data.buffer);
      // eslint-disable-next-line
      item.data = content.toString();
    }
  });
  ctx.body = {
    list: docs,
  };
}

export async function cover(ctx) {
  const no = Joi.attempt(
    ctx.params.no,
    Joi.number()
      .integer()
      .min(0)
      .max(10000)
      .required(),
  );
  const doc = await coverService
    .findOne({
      book: no,
    })
    .lean();
  if (!doc) {
    // TODO 默认封面
    ctx.body = null;
    return;
  }
  ctx.setCache('24h', '5m');
  ctx.set('Content-Type', 'image/jpeg');
  ctx.body = doc.data.buffer;
}
