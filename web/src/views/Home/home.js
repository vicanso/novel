import { mapActions, mapState } from "vuex";
import Banner from "@/components/Banner";
import BookView from "@/components/BookView";
import { routeLogin, routeRegister } from "@/routes";

const functions = {
  shelf: 'shelf',
  hot: 'hot',
  gallery: 'gallery',
  find: 'find',
};

export default {
  name: 'home',
  components: {
    Banner,
    BookView,
  },
  computed: {
    ...mapState({
      books: ({ book }) => book.list,
      bookCount: ({ book }) => book.count,
      bookCategories: ({ book }) => {
        if (!book || !book.categories) {
          return null;
        }
        const {
          categories,
        } = book;
        const result = Object.keys(categories);
        return result.sort((k1, k2) => {
          return categories[k2] - categories[k1];
        });
      },
      bookTodayRecommend: ({ book }) => book.todayRecommend,
      bookLatestPopu: ({ book }) => book.latestPopu,
      bookSearchResult: ({ book }) => book.searchResult,
      userInfo: ({ user }) => user.info,
    })
  },
  data() {
    return {
      banners: null,
      functions,
      currentNav: functions.hot,
      currentCatgory: 0,
      navigation: [
        {
          id: functions.shelf,
          name: '书架',
          cls: 'icon-all',
        },
        {
          id: functions.hot,
          name: '精选',
          cls: 'icon-creditlevel',
        },
        {
          id: functions.gallery,
          name: '书库',
          cls: 'icon-viewgallery',
        },
        {
          id: functions.find,
          name: '发现',
          cls: 'icon-originalimage',
        }
      ],
      field: [
        "id",
        "name",
        "author",
        "brief",
        "cover",
        "wordCount",
      ].join(","),
      hotFields: [
        "id",
        "name",
        "author",
        "brief",
        "cover",
        "wordCount",
        "category",
      ].join(","),
      order: "-updatedAt",
      offset: 0,
      limit: 10,
      loadDone: false,
      loading: false,
      keyword: "",
      searchBooks: null,
    };
  },
  methods: {
    ...mapActions([
      "bookList",
      "bookListCategory",
      "bookListTodayRecommend",
      "bookListLatestPopu",
      "bookCacheRemove",
      "bookUserAction",
      "bookClearSearchResult",
      "bookSearch",
    ]),
    activeNav({id}) {
      if (id === functions.find) {
        this.keyword = "";
      }
      this.currentNav = id;
    },
    reset() {
      this.bookCacheRemove();
      this.loadDone = false;
      this.offset = 0;
    },
    async fetch() {
      const { field, order, offset, limit } = this;
      this.loading = true;
      try {
        const params = {
          field,
          order,
          offset,
          limit,
          status: 2,
        };
        const category = this.bookCategories[this.currentCatgory];
        if (!category) {
          throw new Error('获取失败分类');
        }
        params.category = category;
        await this.bookList(params);
        this.offset = offset + limit;
        if (this.books.length >= this.bookCount) {
          this.loadDone = true;
        }
      } catch (err) {
        this.xError(err);
      } finally {
        this.loading = false;
      }
    },
    async changeCatgeory(index) {
      this.currentCatgory = index;
      this.reset();
      await this.fetch(); 
    },
    initLoadmoreEvent() {
      const {
        loadingMore,
      } = this.$refs;
      const io = new IntersectionObserver(entries => {
        if (this.loading) {
          return;
        }
        const target = entries[0];
        // 在元素可见时加载图标，并做diconnect
        if (target.isIntersecting) {
          this.fetch();
        }
      });
      io.observe(loadingMore);
    },
    async listTodayRecommend() {
      const {
        order,
        hotFields,
      } = this;
      try {
        await this.bookListTodayRecommend({
          limit: 3,
          order,
          field: hotFields,
        });
      } catch (err) {
        this.xError(err);
      }
    },
    async listLatestPopu() {
      const {
        hotFields,
      } = this;
      try {
        await this.bookListLatestPopu({
        limit: 5,
        order: "latestViewCount",
        field: hotFields,
        });
      } catch (err) {
        this.xError(err);
      }
    },
    login() {
      this.$router.push({
        name: routeLogin,
      });
    },
    register() {
      this.$router.push({
        name: routeRegister,
      });
    }
  },
  watch: {
    async keyword(v) {
      if (!v) {
        this.bookClearSearchResult();
        return;
      }
      try {
        this.bookSearch({
          field: 'name,author,id',
          limit: 5,
          keyword: v,
          order: this.order,
        });
      } catch (err) {
        this.xError(err);
      }
    },
  },
  async mounted() {
    const close = this.xLoading();
    try {
      await this.bookListCategory();
      this.fetch();
      this.initLoadmoreEvent();
      this.banners = [
        '01CSRYT09Z31MXY7GG640T9PR2',
        '01CSRYSFXRB7T6W7NPMM2MF478'
      ];
    } catch (err) {
      this.xError(err);
    } finally {
      close();
    }
  },
};