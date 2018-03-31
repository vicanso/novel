import request from 'axios';
import {
  BOOKS,
  BOOKS_DETAIL,
  BOOKS_INFO,
  BOOKS_SOURCES,
} from '../../urls';

const state = {};

const mutations = {};

// 查询书籍
const bookList = async (tmp, query) => {
  const res = await request.get(BOOKS, {
    params: query,
  });
  return res.data;
};

// 更新书籍
const bookUpdate = async (tmp, {no, data}) => {
  const res = await request.patch(BOOKS_DETAIL.replace(':no', no), data);
  return res.data;
};

// 更新书籍信息
const bookUpdateInfo = async (tmp, no) => {
  const res = await request.patch(BOOKS_INFO.replace(':no', no));
  return res.data;
};

// 增加书籍来源
const bookAddSource = async (tmp, {source, name, author, id}) => {
  await request.post(BOOKS_SOURCES, {
    source,
    author,
    name,
    id,
  });
};

export const actions = {
  bookList,
  bookUpdate,
  bookUpdateInfo,
  bookAddSource,
};

export default {
  state,
  mutations,
};