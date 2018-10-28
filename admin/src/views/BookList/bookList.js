import { mapActions, mapState } from "vuex";
import { getListBookPageSize, saveListBookPageSize } from "@/helpers/storage";
import { find } from "lodash-es";

const allCategory = "全部类别";
export default {
  name: "book-list",
  data() {
    return {
      currentBooks: null,
      loading: true,
      tableHeight: 0,
      field: [
        "id",
        "name",
        "author",
        "brief",
        "status",
        "cover",
        "updatedAt",
        "category",
        "sourceCover"
      ].join(","),
      filters: {},
      currentPage: 1,
      order: "-updatedAt",
      offset: 0,
      limit: getListBookPageSize(),
      currentUpdateBoook: null
    };
  },
  computed: {
    ...mapState({
      books: ({ book }) => book.list,
      bookCount: ({ book }) => book.count,
      bookStatusList: ({ book }) => book.statusList,
      bookCategories: ({ book }) => book.categories,
      filterCategories: ({ book }) => {
        if (!book.categories) {
          return [];
        }
        return [allCategory].concat(book.categories);
      }
    })
  },
  methods: {
    ...mapActions([
      "bookList",
      "bookCacheRemove",
      "bookUpdate",
      "bookListCategory",
      "bookGetChapters",
      "bookUpdateCover"
    ]),
    reset() {
      this.bookCacheRemove();
      this.offset = 0;
      this.currentPage = 1;
      this.currentBooks = null;
    },
    handleCurrentChange(page) {
      this.offset = this.limit * (page - 1);
      this.currentPage = page;
      this.fetch();
    },
    handleSizeChange(val) {
      saveListBookPageSize(val);
      this.limit = val;
      this.reset();
      this.fetch();
    },
    search() {
      this.reset();
      this.fetch();
    },
    async fetch() {
      const { field, order, offset, limit } = this;
      const { q, category, status } = this.filters;
      const close = this.xLoading();
      this.loading = true;
      try {
        const params = {
          field,
          order,
          offset,
          limit,
          q
        };
        if (category && category != allCategory) {
          params.category = category;
        }
        if (Number.isInteger(status)) {
          params.status = status;
        }
        await this.bookList(params);
        this.currentBooks = this.books.slice(offset, offset + limit);
      } catch (err) {
        this.xError(err);
      } finally {
        this.loading = false;
        close();
      }
    },
    async edit(book) {
      const close = this.xLoading();
      let latestChapter = null;
      try {
        const res = await this.bookGetChapters({
          id: book.id,
          limit: 1,
          offset: 0,
          order: "-index",
          field: "title",
        });
        const {
          chapters,
        } = res.data;
        if (chapters) {
          latestChapter = chapters[0]
        }
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
      this.currentUpdateBoook = Object.assign({
        latestChapter,
      }, book);
    },
    closeUpdate() {
      this.currentUpdateBoook = null;
    },
    async update() {
      const { currentUpdateBoook } = this;
      const { id, statusDesc } = currentUpdateBoook;
      const found = find(this.currentBooks, item => item.id === id);
      if (!found) {
        return;
      }
      const updateData = {};
      const updateFields = ["brief", "name", "sourceCover"];
      updateFields.forEach(k => {
        const v = currentUpdateBoook[k];
        if (found[k] !== v) {
          updateData[k] = v;
        }
      });
      if (found.statusDesc !== statusDesc) {
        updateData.status = this.bookStatusList.indexOf(statusDesc);
      }
      const currentCatgory = found.category.sort().join(",");
      const newCategory = currentUpdateBoook.category.sort().join(",");
      if (currentCatgory != newCategory) {
        updateData.category = newCategory;
      }
      const close = this.xLoading();
      try {
        await this.bookUpdate({
          id,
          update: updateData
        });
        this.currentUpdateBoook = null;
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    },
    async updateCover() {
      const { id } = this.currentUpdateBoook;
      const close = this.xLoading();
      try {
        await this.bookUpdateCover({
          id
        });
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    }
  },
  async mounted() {
    const paginationHeight = 60;
    const filterHeight = 50;
    this.tableHeight =
      this.$refs.bookList.clientHeight - paginationHeight - filterHeight;
  },
  beforeMount() {
    this.fetch();
  }
};