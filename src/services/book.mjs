import spider from 'novel-spider';
import request from 'superagent';
import Promise from 'bluebird';
import zlib from 'zlib';
import util from 'util';
import _ from 'lodash';

import genService from './gen';
import originService from './origin';
import coverService from './cover';
import chapterService from './chapter';
import {getBookNo} from './inc';
import redis, {lock} from '../helpers/redis';
import errors from '../errors';

const bookService = genService('Book');
const gzip = util.promisify(zlib.gzip);

// 获取对应的spider
async function getSpider(author, name) {
  const doc = await originService.findOne({
    name,
    author,
  });
  if (!doc || doc.source !== 'biquge') {
    return null;
  }
  const id = doc.sourceId;
  return new spider.BiQuGe(id);
}

export default bookService;

// 添加一本新书，如果已有，则返回
export async function addBook(author, name) {
  let doc = await bookService.findOne({
    author,
    name,
  });
  if (doc) {
    return doc;
  }
  const novel = await getSpider(author, name);
  if (!novel) {
    throw errors.get('book.sourceNotFound');
  }
  const info = await novel.getInfos();
  if (!info || info.name !== name || info.author !== author) {
    return null;
  }
  const no = await getBookNo();
  const data = {
    no,
    name,
    author,
    brief: info.brief,
  };
  doc = await bookService.add(data);
  const count = await coverService.count({
    book: no,
  });
  if (count === 0) {
    const res = await request.get(info.img);
    await coverService.add({
      book: no,
      data: res.body,
    });
  }
  return doc;
}

// 更新书本最新章节
export async function updateChapters(author, name) {
  const doc = await bookService.findOne(
    {
      author,
      name,
    },
    'no',
  );
  if (!doc) {
    return 0;
  }
  const bookNo = doc.no;
  const chapterCount = await chapterService.count({book: bookNo});
  const novel = await getSpider(author, name);
  const chapters = await novel.getChapters();
  const updateChapter = async (chapter, i) => {
    const index = i + chapterCount;
    const {title, content} = await novel.getChapter(index);
    const gzipContent = await gzip(Buffer.from(content));
    await chapterService.add({
      book: bookNo,
      no: index,
      wordCount: content.length,
      title,
      data: gzipContent,
    });
  };
  await Promise.mapSeries(chapters.slice(chapterCount), updateChapter);
  return chapters.length - chapterCount;
}

// 更新书相关信息（最近更新时间，字数，章节数等）
export async function updateInfo(author, name) {
  const doc = await bookService.findOne({
    author,
    name,
  });
  if (!doc) {
    return;
  }
  const docs = await chapterService
    .find({
      book: doc.no,
    })
    .select('wordCount no title updatedAt')
    .lean();
  if (docs.length === 0 || doc.chapterCount === docs.length) {
    return;
  }
  let wordCount = 0;
  _.forEach(docs, item => {
    wordCount += item.wordCount;
  });
  doc.latestChapter = _.pick(_.last(docs), [
    'title',
    'wordCount',
    'updatedAt',
    'no',
  ]);
  if (!doc.end) {
    doc.end = false;
  }
  doc.wordCount = wordCount;
  doc.chapterCount = docs.length;
  await doc.save();
}

// 更新所有书籍
export async function updateAll() {
  const docs = await bookService
    .find({
      end: {
        $ne: true,
      },
    })
    .select('author name no')
    .lean();
  let count = 0;
  const ttl = 300;
  await Promise.mapSeries(docs, async doc => {
    const key = `update-all-books-${doc.no}`;
    const locked = await lock(key, ttl);
    // 如果出错或者setnx不成功（有其它实例已在更新）
    if (!locked) {
      return;
    }
    console.info(`the book(${doc.no}) will be updated`);
    const {author, name} = doc;
    await updateChapters(author, name);
    await updateInfo(author, name);
    count += 1;
  });
  console.info(`update ${count} books is finished`);
}

// 获取书籍分类汇总
export async function getCategories() {
  const key = 'book-categories';
  const data = await redis.get(key);
  setImmediate(async () => {
    // 控制最多每10分钟更新一次
    const locked = await lock('update-categories', 60 * 10);
    if (!locked) {
      return;
    }
    const cursor = bookService
      .find({})
      .select('category')
      .lean()
      .cursor();
    const result = {};
    cursor.on('data', doc => {
      _.forEach(doc.category, category => {
        if (!result[category]) {
          result[category] = 0;
        }
        result[category] += 1;
      });
    });
    cursor.on('end', async () => {
      try {
        const list = [];
        _.forEach(result, (v, k) => {
          list.push({
            name: k,
            count: v,
          });
        });
        await redis.set(
          key,
          JSON.stringify({
            createdAt: new Date().toISOString(),
            categories: _.sortBy(list, item => -item.count),
          }),
        );
      } catch (err) {
        console.error(`save categories to redis fail, ${err.message}`);
      }
    });
    cursor.on('error', err => {
      console.error(`get all book categories fail, ${err.message}`);
    });
  });
  if (!data) {
    return null;
  }
  return JSON.parse(data);
}

// 更新书籍分类
export async function updateAllCategory() {
  // 对已经完结的增加 完本 分类
  const docs = await bookService
    .find({
      end: true,
    })
    .select('author name category');
  await Promise.mapSeries(docs, async doc => {
    if (!_.includes(doc.category, '完本')) {
      doc.category.push('完本');
      await doc.save();
    }
  });
}
