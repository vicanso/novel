import request from "axios";

import { BOOKS } from "@/urls";
import { BOOK_LIST } from "@/store/types";

import { debug, formatDate } from "@/helpers/util";

const statusList = ["待审核", "已拒绝", "已通过"];

const state = {
  book: {
    list: [],
    count: 0,
    statusList,
  }
};

const bookList = async ({ commit }, { field, order, offset, limit }) => {
  const params = {
    field,
    order,
    offset,
    limit
  };
  const { list } = state.book;
  if (list[offset]) {
    return;
  }
  debug(params);
  const res = await request.get(BOOKS, {
    params
  });
  debug(res);
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

const actions = {
  bookList,
  bookCacheRemove
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
      const status = item.status || 0;
      item.statusDesc = statusList[status];

      stateData.list[offset + i] = item;
    });
    if (count >= 0) {
      stateData.count = count;
    }
  }
};

export default {
  actions,
  state,
  mutations
};
