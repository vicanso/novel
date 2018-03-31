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
      categoryOptions: ['玄幻', '都市', '仙侠', '科幻', '游戏', '历史'],
      sourceOptions: [
        {
          label: '笔趣阁',
          value: 'biquge',
        },
      ],
      sourceForm: null,
    };
  },
  methods: {
    ...mapActions(['bookList', 'bookUpdate', 'bookUpdateInfo', 'bookAddSource']),
    changePage(page) {
      this.page = page - 1;
      this.loadBooks();
    },
    edit(no) {
      this.form = _.clone(_.find(this.books, item => item.no === no));
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
      const found =_.find(books, item => item.no === no)
      const data = diff(found, form);
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
      } = this;
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
