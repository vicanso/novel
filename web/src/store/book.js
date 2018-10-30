import request from "axios";
import {
  BOOKS,
  BOOKS_RECOMMEND_BY_ID,
  BOOKS_CATEGORIES,
  BOOKS_CHAPTERS,
  BOOKS_USER_ACTIONS,
  BOOKS_DETAIL
} from "@/urls";
import {
  BOOK_LIST,
  BOOK_CATEGORY,
  BOOK_LIST_TODAY_RECOMMEND,
  BOOK_SEARCH_RESULT,
  BOOK_LIST_LATEST_POPU
} from "@/store/types";

import {
  ChapterCache,
  BookReadInfo,
  getStoreChapterIndexList
} from "@/helpers/storage";

import { formatDate } from "@/helpers/util";

const todayHotCategory = "今日必读";
const statusPassed = 2;
var currentKeyword = "";

const state = {
  book: {
    detail: null,
    list: null,
    count: 0,
    categories: null,
    todayRecommend: null,
    latestPopu: null,
    searchResult: null
  }
};

const bookList = async (
  { commit },
  { field, order, offset, limit, q, category, status }
) => {
  const params = {
    field,
    order,
    offset,
    limit
  };
  if (q) {
    params.q = q;
  }
  if (category) {
    params.category = category;
  }
  if (Number.isInteger(status)) {
    params.status = status;
  }
  const { list } = state.book;
  if (list && list[offset]) {
    return;
  }
  const res = await request.get(BOOKS, {
    params
  });
  commit(
    BOOK_LIST,
    Object.assign(
      {
        offset
      },
      res.data
    )
  );
};

const bookCacheRemove = async ({ commit }) => {
  commit(BOOK_LIST, null);
};

const bookListCategory = async ({ commit }) => {
  const res = await request.get(BOOKS_CATEGORIES);
  commit(BOOK_CATEGORY, res.data);
  return res;
};

const bookListTodayRecommend = async ({ commit }, { limit, field, order }) => {
  const params = {
    field,
    offset: 0,
    limit,
    order,
    category: todayHotCategory,
    status: statusPassed
  };
  const res = await request.get(BOOKS, {
    params
  });
  commit(BOOK_LIST_TODAY_RECOMMEND, res.data.books);
  return res;
};

const bookListLatestPopu = async ({ commit }, { limit, field, order }) => {
  const params = {
    field,
    offset: 0,
    limit,
    order,
    status: statusPassed
  };
  const res = await request.get(BOOKS, {
    params
  });
  commit(BOOK_LIST_LATEST_POPU, res.data.books);
  return res;
};

const bookUserAction = async (tmp, { id, type }) => {
  const url = BOOKS_USER_ACTIONS.replace(":id", id);
  const res = await request.post(url, {
    type
  });
  return res;
};

const bookSearch = async ({ commit }, { keyword, field, limit, order }) => {
  currentKeyword = keyword;
  const params = {
    field,
    offset: 0,
    limit,
    q: keyword,
    order,
    status: statusPassed
  };
  const res = await request.get(BOOKS, {
    params
  });
  if (keyword === currentKeyword) {
    commit(BOOK_SEARCH_RESULT, res.data.books);
  }
  return res;
};

const bookClearSearchResult = async ({ commit }) => {
  commit(BOOK_SEARCH_RESULT, null);
};

const bookGetDetail = async (tmp, { id }) => {
  const res = await request.get(BOOKS_DETAIL.replace(":id", id));
  return res;
};

// bookGetRecommend get recommend
const bookGetRecommend = async (tmp, { id, limit, field, order }) => {
  const url = BOOKS_RECOMMEND_BY_ID.replace(":id", id);
  const res = await request.get(url, {
    params: {
      status: statusPassed,
      offset: 0,
      limit,
      field,
      order
    }
  });
  return res;
};

const bookGetChapters = async (tmp, { id, limit, offset, field, order }) => {
  const url = BOOKS_CHAPTERS.replace(":id", id);
  const res = await request.get(url, {
    params: {
      limit,
      offset,
      field,
      order
    }
  });
  return res;
};

const getChapters = async (id, offset, limit) => {
  const url = BOOKS_CHAPTERS.replace(":id", id);
  return request.get(url, {
    params: {
      limit,
      offset,
      field: "title,content",
      order: "index"
    }
  });
};

const bookGetChapterContent = async (tmp, { id, no }) => {
  const c = new ChapterCache(id);
  const data = await c.get(no);
  if (data) {
    return data;
  }
  const limit = 10;
  let offset = Math.floor(no / limit) * limit;
  const res = await getChapters(id, offset, limit);
  const { chapters } = res.data;
  chapters.forEach((item, index) => {
    c.add(offset + index, item);
  });
  return chapters[no - offset];
};

const bookDownload = async (tmp, { id, max }) => {
  const limit = 20;
  const arr = [];
  for (let index = 0; index < Math.ceil(max / limit); index++) {
    arr.push(index);
  }
  const storeIndexList = await getStoreChapterIndexList(id);
  const dict = {};
  storeIndexList.forEach(v => {
    dict[v] = true;
  });
  const c = new ChapterCache(id);
  return Promise.map(
    arr,
    async i => {
      const offset = i * limit;
      let found = false;
      for (let index = 0; index < limit; index++) {
        if (found) {
          break;
        }
        const v = Math.min(index + offset, max - 1);
        // 如果发现有章节未下载
        if (!dict[v]) {
          found = true;
        }
      }
      // 如果未发现未下载章节，则无需下载
      if (!found) {
        return;
      }
      const res = await getChapters(id, offset, limit);
      const { chapters } = res.data;
      chapters.forEach((item, index) => {
        c.add(offset + index, item);
      });
    },
    {
      concurrency: 1
    }
  );
};

// 获取缓存的章节序号
const bookGetStoreChapterIndexes = async (tmp, { id }) => {
  return await getStoreChapterIndexList(id);
};

// bookGetReadInfo 获取当前阅读信息（阅读至第几章，开始阅读时间，最新阅读时间）
const bookGetReadInfo = async (tmp, { id }) => {
  const b = new BookReadInfo(id);
  return await b.get();
};

// bookUpdateReadInfo 更新当前阅读信息
const bookUpdateReadInfo = async (tmp, { id, no, page }) => {
  const b = new BookReadInfo(id);
  await b.update(no, page);
};

const actions = {
  bookGetDetail,
  bookList,
  bookCacheRemove,
  bookListTodayRecommend,
  bookListCategory,
  bookListLatestPopu,
  bookSearch,
  bookClearSearchResult,
  bookGetRecommend,
  bookGetChapters,
  bookGetChapterContent,
  bookDownload,
  bookGetStoreChapterIndexes,
  bookGetReadInfo,
  bookUpdateReadInfo,
  bookUserAction
};

const mutations = {
  [BOOK_LIST](state, data) {
    const stateData = state.book;
    // clear cache
    if (!data) {
      stateData.list = null;
      stateData.count = 0;
      return;
    }
    if (!stateData.list) {
      stateData.list = [];
    }
    const { books, count, offset } = data;
    books.forEach(function(item, i) {
      if (item.updatedAt) {
        item.updatedAt = formatDate(item.updatedAt);
      }
      stateData.list[offset + i] = item;
    });
    if (count >= 0) {
      stateData.count = count;
    }
  },
  [BOOK_CATEGORY](state, { categories }) {
    state.book.categories = categories;
  },
  [BOOK_LIST_TODAY_RECOMMEND](state, data) {
    state.book.todayRecommend = data;
  },
  [BOOK_LIST_LATEST_POPU](state, data) {
    state.book.latestPopu = data;
  },
  [BOOK_SEARCH_RESULT](state, data) {
    state.book.searchResult = data;
  }
};

export default {
  actions,
  state,
  mutations
};
