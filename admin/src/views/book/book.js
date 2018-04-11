import {mapActions} from 'vuex';
import _ from 'lodash';

import {
  getDate,
  diff,
} from '../../helpers/util';


export default {
  data() {
    return {
      page: 0,
      size: 10,
      count: -1,
      books: null,
      form: null,
      mode: -1,
      categoryOptions: ['玄幻', '都市', '仙侠', '科幻', '游戏', '历史', '灵异', '军事', '体育', '二次元'],
      statusList: [
        '所有',
        '完结',
        '未完结',
      ],
      sourceOptions: [
        {
          label: '笔趣阁',
          value: 'biquge',
        },
      ],
      sourceForm: null,
      searchForm: {},
    };
  },
  methods: {
    ...mapActions(['bookList', 'bookUpdate', 'bookUpdateInfo', 'bookAddSource']),
    search() {
      this.page = 0;
      this.loadBooks();
    },
    changePage(page) {
      this.page = page - 1;
      this.loadBooks();
    },
    edit(no) {
      this.form = _.cloneDeep(_.find(this.books, item => item.no === no));
      this.mode = 1;
    },
    addSource() {
      this.sourceForm = {};
      this.mode = 2;
    },
    async submitSource() {
      const {source, name, author, id} = this.sourceForm;
      if (!source || !name || !author || !id) {
        this.$error('必填字段不能为空');
        return;
      }
      const close = this.$loading();
      try {
        await this.bookAddSource({
          source,
          name,
          author,
          id,
        });
        this.$message('书籍来源添加成功');
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
    async update(no) {
      const close = this.$loading();
      try {
        await this.bookUpdateInfo(no);
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
    async submit() {
      const {
        form,
        books,
      } = this;
      const {
        no,
      } = form;
      const keys = [
        'end',
        'category',
        'brief'
      ];
      const found =_.find(books, item => item.no === no)
      const data = diff(found, form, keys);
      const close = this.$loading();
      try {
        await this.bookUpdate({
          no,
          data,
        });
        _.extend(found, data);
        this.mode = 0;
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
    async loadBooks() {
      const {
        page,
        size,
        searchForm,
      } = this;
      const {
        keyword,
        status,
      } = searchForm;
      const close = this.$loading();
      try {
        const query = {
          skip: page * size,
          limit: size,
          sort: '-updatedAt',
        };
        if (page === 0) {
          query.count = true;
        }
        if (keyword) {
          query.keyword = keyword;
        }
        if (status === '完结') {
          query.end = true;
        } else if (status === '未完结') {
          query.end = false;
        }
        const data = await this.bookList(query);
        if (page === 0) {
          this.count = data.count;
        }
        this.books = _.map(data.books, (item) => {
          item.date = getDate(item.updatedAt);
          return item;
        });
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
  },
  beforeMount() {
    this.loadBooks().then(() => {
      this.mode = 0;
    });
  },
};
