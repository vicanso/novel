import Joi from 'joi';
import Promise from 'bluebird';
import zlib from 'zlib';
import util from 'util';

import orginService from '../services/origin';
import bookService, {
  addBook,
  updateChapters,
  updateInfo,
  getCategories,
} from '../services/book';
import chapterService from '../services/chapter';
import coverService from '../services/cover';
import {lock} from '../helpers/redis';

const gunzip = util.promisify(zlib.gunzip);

const schema = {
  name: () =>
    Joi.string()
      .trim()
      .max(30),
  author: () =>
    Joi.string()
      .trim()
      .max(20),
  source: () =>
    Joi.string()
      .trim()
      .valid(['biquge']),
  id: () =>
    Joi.string()
      .trim()
      .max(10),
  no: () =>
    Joi.number()
      .integer()
      .min(0)
      .max(10000),
  category: () =>
    Joi.string()
      .trim()
      .max(10),
  sort: () =>
    Joi.string()
      .trim()
      .max(30),
};

// 增加来源
export async function addSource(ctx) {
  const {name, author, source, id} = Joi.validate(ctx.request.body, {
    name: schema.name().required(),
    author: schema.author().required(),
    source: schema.source().required(),
    id: schema.id().required(),
  });
  const doc = await orginService.add({
    name,
    author,
    source,
    sourceId: id,
  });
  await addBook(author, name);
  ctx.status = 201;
  ctx.body = doc;
}

// 获取书籍列表
export async function list(ctx) {
  const {skip, limit, count, fields, keyword, category, sort} = Joi.validate(
    ctx.query,
    {
      skip: Joi.number()
        .integer()
        .default(0),
      limit: Joi.number()
        .integer()
        .min(1)
        .max(20)
        .default(10),
      fields: Joi.string()
        .trim()
        .max(100),
      keyword: Joi.string()
        .trim()
        .max(30),
      category: schema.category(),
      sort: schema.sort(),
      count: Joi.boolean(),
    },
  );
  const conditions = {};
  if (keyword) {
    conditions.keyword = new RegExp(keyword);
  }
  if (category) {
    conditions.category = category;
  }
  const data = {};
  if (count) {
    data.count = await bookService.count(conditions);
  }
  data.list = await bookService
    .find(conditions)
    .skip(skip)
    .limit(limit)
    .sort(sort)
    .select(fields)
    .lean();
  ctx.setCache(60);
  ctx.body = data;
}

// 增加书籍
export async function add(ctx) {
  const {name, author} = Joi.validate(ctx.request.body, {
    name: schema.name().required(),
    author: schema.author().required(),
  });
  const doc = await addBook(author, name);
  ctx.status = 201;
  ctx.body = doc;
}

// 获取书籍信息
export async function get(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
  const doc = await bookService.findOne({
    no,
  });
  ctx.setCache('5m');
  ctx.body = doc;
}

// 更新书籍章节信息
export async function updateBookInfo(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
  const locked = await lock(`updte-book-${no}`, 300);
  const doc = await bookService
    .findOne({
      no,
    })
    .select('author name')
    .lean();
  const {author, name} = doc;
  if (locked) {
    updateChapters(author, name)
      .then(() => updateInfo(author, name))
      .then(() => {
        console.info(`update book(${no}) success`);
      })
      .catch(err => {
        console.error(`update book(${no}) fail, ${err.message}`);
      });
  }
  ctx.body = null;
}

// 更新书籍信息
export async function update(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
  const {brief, end, category} = Joi.validate(ctx.request.body, {
    brief: Joi.string()
      .trim()
      .max(500),
    end: Joi.boolean(),
    category: Joi.array().items(Joi.string().trim()),
  });
  const doc = await bookService.findOne({
    no,
  });
  if (brief) {
    doc.brief = brief;
  }
  if (category) {
    doc.category = category;
  }
  doc.end = end;
  await doc.save();
  ctx.body = doc;
}

// 列出章节信息
export async function listChapter(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
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

// 获取封面信息
export async function getCover(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
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
  ctx.setCache('1w', '5m');
  ctx.set('Content-Type', 'image/jpeg');
  ctx.body = doc.data.buffer;
}

// 获取分类信息
export async function categoriesList(ctx) {
  const data = await getCategories();
  ctx.setCache('5m');
  ctx.body = data;
}
