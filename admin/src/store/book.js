import request from "axios";
import { find } from "lodash-es";

import {
  BOOKS,
  BOOKS_UPDATE_INFO,
  BOOKS_UPDATE_COVEER,
  BOOKS_CHAPTERS
} from "@/urls";
import { BOOK_LIST, BOOK_UPDATE, BOOK_UPDATE_COVER } from "@/store/types";

import { formatDate } from "@/helpers/util";

const statusList = ["待审核", "已拒绝", "已通过"];

const state = {
  book: {
    list: [],
    count: 0,
    statusList,
    categories: [
      "今日必读",
      "玄幻奇幻",
      "女频频道",
      "都市言情",
      "武侠仙侠",
      "历史军事",
      "科幻灵异",
      "网游竞技",
      "完结"
    ]
  }
};

function getStatusDesc(status) {
  if (!status) {
    return statusList[0];
  }
  return statusList[status];
}

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
  if (list[offset]) {
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

const bookUpdate = async ({ commit }, { id, update }) => {
  const url = BOOKS_UPDATE_INFO.replace(":id", id);
  await request.patch(url, update);
  const data = Object.assign({}, update);
  if (update.category) {
    data.category = update.category.split(",");
  }
  commit(BOOK_UPDATE, {
    id,
    update: data
  });
};

const bookCacheRemove = async ({ commit }) => {
  commit(BOOK_LIST, null);
};

const bookUpdateCover = async ({ commit }, { id }) => {
  const res = await request.patch(BOOKS_UPDATE_COVEER.replace(":id", id));
  commit(
    BOOK_UPDATE_COVER,
    Object.assign(
      {
        id
      },
      res.data
    )
  );
};

const bookGetChapters = async (tmp, { id, offset, limit, order, field }) => {
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

const actions = {
  bookList,
  bookCacheRemove,
  bookUpdate,
  bookGetChapters,
  bookUpdateCover
};

const mutations = {
  [BOOK_LIST](state, data) {
    const stateData = state.book;
    // clear cache
    if (!data) {
      stateData.list = [];
      stateData.count = 0;
      return;
    }
    const { books, count, offset } = data;
    books.forEach(function(item, i) {
      if (item.updatedAt) {
        item.updatedAt = formatDate(item.updatedAt);
      }
      item.statusDesc = getStatusDesc(item.status);

      stateData.list[offset + i] = item;
    });
    if (count >= 0) {
      stateData.count = count;
    }
  },
  [BOOK_UPDATE](state, { id, update }) {
    const found = find(state.book.list, item => item.id === id);
    if (found) {
      Object.assign(found, update);
      found.statusDesc = getStatusDesc(update.status);
    }
  },
  [BOOK_UPDATE_COVER](state, { id, cover }) {
    const found = find(state.book.list, item => item.id === id);
    if (found) {
      found.cover = cover;
    }
  }
};

export default {
  actions,
  state,
  mutations
};
