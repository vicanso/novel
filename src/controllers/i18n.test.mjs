import shortid from 'shortid';

import '../helpers/joi-extend';

import {add, update, list, listCategory, listCategoryByLang, get} from './i18n';
import i18nService from '../services/i18n';
import {createContext} from '../helpers/utils';

const randomId = shortid();
const account = 'vicanso';

afterAll(async () => {
  await i18nService.Model.remove({
    name: randomId,
  });
});

test('add i18n', async () => {
  const data = {
    name: randomId,
    category: randomId,
    en: 'title',
    zh: '标题',
  };
  const ctx = createContext('/i18ns');
  ctx.session = {
    user: {
      account,
    },
  };
  ctx.request.body = data;
  await add(ctx);
  expect(ctx.status).toBe(201);
  Object.keys(data).forEach(k => {
    expect(data[k]).toBe(ctx.body[k]);
  });
  expect(ctx.body.creator).toBe(account);
});

test('update i18n', async () => {
  let doc = await i18nService
    .findOne({
      name: randomId,
    })
    .lean();
  const en = 'title-new';
  // eslint-disable-next-line
  const id = doc._id.toString();
  const ctx = createContext(`/i18ns/${id}`);
  ctx.params.id = id;
  ctx.request.body = {
    en,
  };
  await update(ctx);
  expect(ctx.status).toBe(204);
  doc = await i18nService
    .findOne({
      name: randomId,
    })
    .lean();
  expect(doc.en).toBe(en);
});

test('list i18n', async () => {
  const ctx = createContext(`/i18ns?category=${randomId}&count=true`);
  await list(ctx);
  expect(ctx.status).toBe(200);
  expect(ctx.body.count).toBe(1);
  expect(ctx.body.list[0].name).toBe(randomId);
});

test('list i18n categories', async () => {
  const ctx = createContext('/i18ns/categories');
  await listCategory(ctx);
  expect(ctx.status).toBe(200);
  expect(ctx.body.list.length).toBeGreaterThan(0);
});

test('list i18n category by lang', async () => {
  const category = randomId;
  const name = randomId;
  const ctx = createContext(`/i18ns/category?category=${randomId}`);
  ctx.state.lang = 'en';
  await listCategoryByLang(ctx);
  expect(ctx.status).toBe(200);
  expect(ctx.body[category][name]).toBe('title-new');
});

test('get i18n', async () => {
  const doc = await i18nService
    .findOne({
      name: randomId,
    })
    .lean();
  // eslint-disable-next-line
  const id = doc._id.toString();
  const ctx = createContext(`/i18ns/${id}`);
  ctx.params.id = id;
  await get(ctx);
  expect(ctx.status).toBe(200);
  // eslint-disable-next-line
  expect(ctx.body.name).toBe(randomId);
});
