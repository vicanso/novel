import Joi from 'joi';
import _ from 'lodash';
import Promise from 'bluebird';
import zlib from 'zlib';
import util from 'util';

import orginService from '../services/origin';
import bookService, {
  addBook,
  updateChapters,
  updateInfo,
  getCategories,
  getSources,
  updateCover,
} from '../services/book';
import chapterService from '../services/chapter';
import coverService from '../services/cover';
import {lock} from '../helpers/redis';
import requestBookService from '../services/request-book';
import * as tinyService from '../services/tiny';

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
      .valid(getSources()),
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
  brief: () =>
    Joi.string()
      .trim()
      .max(500),
  end: () => Joi.boolean(),
  skip: () =>
    Joi.number()
      .integer()
      .default(0),
  limit: () =>
    Joi.number()
      .integer()
      .min(1)
      .max(100)
      .default(10),
  fields: () =>
    Joi.string()
      .trim()
      .max(100),
  keyword: () =>
    Joi.string()
      .trim()
      .max(30),
  imageType: () =>
    Joi.string()
      .valid(['webp', 'jpeg'])
      .default('jpeg'),
  imageQuality: () =>
    Joi.number()
      .integer()
      .min(70)
      .max(100)
      .default(75),
};

// 增加来源
export async function addSource(ctx) {
  const {name, author, source, id} = Joi.validate(ctx.request.body, {
    name: schema.name().required(),
    author: schema.author().required(),
    source: schema.source().required(),
    id: schema.id().required(),
  });
  let doc = await orginService
    .findOne({
      name,
      author,
    })
    .lean();
  if (!doc) {
    doc = await orginService.add({
      name,
      author,
      source,
      sourceId: id,
    });
  }
  await addBook(author, name);
  ctx.status = 201;
  ctx.body = doc;
}

// 获取书籍列表
export async function list(ctx) {
  const {
    skip,
    limit,
    count,
    fields,
    keyword,
    category,
    sort,
    no,
    end,
  } = Joi.validate(ctx.query, {
    skip: schema.skip(),
    limit: schema.limit(),
    fields: schema.fields(),
    keyword: schema.keyword(),
    no: Joi.string()
      .trim()
      .max(300),
    category: schema.category(),
    sort: schema.sort(),
    end: schema.end(),
    count: Joi.boolean(),
  });
  const conditions = {};
  if (keyword) {
    conditions.keyword = new RegExp(keyword);
  }
  if (category) {
    conditions.category = category;
  }
  if (no) {
    conditions.no = {
      $in: _.map(no.split(','), v => Number.parseInt(v, 10)),
    };
  }
  if (!_.isUndefined(end)) {
    if (end) {
      conditions.end = true;
    } else {
      conditions.end = {
        $ne: true,
      };
    }
  }
  const data = {};
  if (count) {
    data.count = await bookService.count(conditions);
  }
  data.books = await bookService
    .find(conditions)
    .skip(skip)
    .limit(limit)
    .sort(sort)
    .select(fields)
    .lean();
  ctx.setCache('10s');
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
    brief: schema.brief(),
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
  if (!_.isUndefined(end)) {
    doc.end = end;
  }
  await doc.save();
  ctx.body = doc;
}

// 列出章节信息
export async function listChapter(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
  const {skip, limit, fields, sort} = Joi.validate(ctx.query, {
    skip: schema.skip(),
    limit: schema.limit(),
    fields: schema.fields(),
    sort: schema.sort().default('updatedAt'),
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
  // 如果不等于，表示已到最后章节
  // 因此缓存时间设置较短
  if (limit === docs.length) {
    ctx.setCache('30m', '5m');
  } else {
    ctx.setCache('5m');
  }
  ctx.body = {
    chapters: docs,
  };
}

// 获取封面信息
export async function getCover(ctx) {
  const {type, quality} = Joi.validate(ctx.query, {
    type: schema.imageType(),
    quality: schema.imageQuality(),
  });
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
  let buf = doc.data.buffer;
  let contentType = type;
  try {
    let fn = tinyService.toJpeg;
    if (type === 'webp') {
      fn = tinyService.toWebp;
    }
    const res = await fn(buf, quality);
    buf = res.data;
  } catch (err) {
    contentType = 'jpeg';
    console.error(`conver cover to ${type} fail, ${err.message}`);
  }
  ctx.setCache('1w', '5m');
  ctx.set('Content-Type', `image/${contentType}`);
  ctx.body = buf;
}

// 更新封面
export async function coverUpdate(ctx) {
  const no = Joi.attempt(ctx.params.no, schema.no().required());
  const result = await updateCover(no);
  ctx.body = {
    success: result,
  };
}

// 获取分类信息
export async function categoriesList(ctx) {
  const data = await getCategories();
  ctx.setCache('5m');
  ctx.body = data;
}

// 请求增加书籍
export async function requestBook(ctx) {
  const data = Joi.validate(ctx.request.body, {
    name: schema.name().required(),
    author: schema.author().required(),
  });
  const doc = await requestBookService.findOne(data);
  if (!doc) {
    await requestBookService.add(data);
  }
  ctx.status = 201;
}

// 获取推荐书籍
export async function getRecommendations(ctx) {
  const {limit} = Joi.validate(ctx.query, {
    limit: schema.limit().default(3),
  });
  const no = Joi.attempt(ctx.params.no, schema.no().required());
  const fields = 'name no author';
  const doc = await bookService
    .findOne({
      no,
    })
    .select('author category')
    .lean();
  // 选择作者相同小说
  let docs = await bookService
    .find({
      author: doc.author,
      no: {
        $ne: no,
      },
    })
    .select(fields)
    .limit(limit)
    .lean();
  const result = [];
  _.forEach(docs, item => {
    if (_.random(0, 10) < 5) {
      result.push(item);
    }
  });

  const addMore = async conditions => {
    docs = await bookService
      .find(conditions)
      .select(fields)
      .limit(limit)
      .lean();
    docs = _.shuffle(docs);
    result.push(...docs.slice(0, limit - result.length));
  };

  // 增加同类型小说
  await addMore({
    no: {
      $ne: no,
    },
    category: _.sample(doc.category),
  });
  // 如果同类型的还不足够，随机选择
  if (result.length < limit) {
    await addMore({
      no: {
        $ne: no,
      },
    });
  }
  ctx.setCache('10m');
  ctx.body = {
    recommendations: result,
  };
}
