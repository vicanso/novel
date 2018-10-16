<template lang="pug">
mixin BookTable
  el-table.table(
    :data="currentBooks"
    :height="tableHeight"
    border
    stripe
  )
    el-table-column(
      prop="updatedAt"
      label="更新时间"
      width="160"
    )
    el-table-column(
      prop="author"
      label="作者"
      width="120"
    )
    el-table-column(
      prop="name"
      label="书名"
      width="160"
    )
    el-table-column(
      prop="brief"
      label="简介"
    )
    el-table-column(
      prop="statusDesc"
      label="状态"
      width="80"
    )
    el-table-column(
      fixed="right"
      label="操作"
      width="60"
    )
      template(
        slot-scope="scope"
      )
        el-button(
          @click="edit(scope.row)"
          type="text"
          size="small"
        ) 编辑
  .pagination: el-pagination(
    layout="sizes, prev, pager, next"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
    :page-sizes="[10, 20, 30, 50]"
    :pageSize="limit" 
    :total="bookCount"
    :current-page="currentPage"
  )

mixin BookUpdate
  .bookUpdate(
    v-if="currentUpdateBoook"
  )
    h5
      a.pullRight.close(
        href='javascript:;'
        @click="closeUpdate()"
      )
        i.el-icon-close
      | 编辑内容
    el-form.form(
      label-width="90px"
    )
      el-form-item(
        label="作者"
      )
        el-input(
          :disabled="true"
          v-model="currentUpdateBoook.author"
        )
      el-form-item(
        label="书名"
      )
        el-input(
          v-model="currentUpdateBoook.name"
        )
      el-form-item(
        label="状态"
      )
        el-select(
          placeholder="请选择状态"
          v-model="currentUpdateBoook.statusDesc"
        )
          el-option(
            v-for="status in bookStatusList"
            :key="status"
            :label="status"
            :value="status"
          )
      el-form-item(
        label="分类"
      )
        el-select(
          style="width:100%"
          placeholder="请选择分类"
          multiple
          v-model="currentUpdateBoook.category"
        )
          el-option(
            v-for="item in bookCategories"
            :key="item"
            :label="item"
            :value="item"
          )
      el-form-item(
        label="原始封面"
      )
        img(
          :src="currentUpdateBoook.sourceCover"
          height="60px"
        ) 
        el-button.mleft10(
          type="text"
          v-if="!currentUpdateBoook.cover"
          @click.native="updateCover"
        ) 更新
      el-form-item(
        label="简介"
      )
        el-input(
          v-model="currentUpdateBoook.brief"
          type="textarea"
          :autosize="{ minRows: 4 }"
        )
      el-form-item
        el-button(
          type="primary"
          style="width:100%"
          @click.native="update"
        ) 保存

mixin BookFilter
  .bookFilter
    el-input(
      v-model="filters.q"
      placeholder="请输入关键字"
      clearable
    )
      el-select.categorySelector(
        slot="prepend"
        placeholder="分类"
        v-model="filters.category"
      )
        el-option(
          v-for="category in filterCategories"
          :key="category"
          :label="category"
          :value="category"
        )
      div(
        slot="append"
      )
        el-radio.status(
          v-model="filters.status"
          v-for="status, index in bookStatusList"
          :key="status"
          :label="index"
        ) {{status}}
        span.divide |
        el-button.search(
          icon="el-icon-search"
          @click.native="search"
        )
.bookList(
  ref="bookList"
)
  +BookFilter
  .tableWrapper(
    v-if="!loading"
  )
    +BookTable
  +BookUpdate

</template>
<style lang="sass">
@import "@/styles/const.sass"
.bookList
  position: fixed 
  top: $MAIN_HEADER_HEIGHT
  left: $MAIN_NAV_WIDTH
  right: 0
  bottom: 0
  padding: 10px
.bookUpdate
  position: absolute
  top: 50%
  left: 50%
  width: 600px
  margin-left: -300px
  margin-top: -300px
  border: $GRAY_BORDER
  background-color: $COLOR_WHITE
  border-radius: 3px
  z-index: 9
  h5
    margin: 0
    padding: 0
    line-height: 3em
    padding-left: 10px
    background-color: $COLOR_BLACK
    color: $COLOR_WHITE
    .close
      color: $COLOR_WHITE
      display: block
      width: 40px
      text-align: center
      font-size: 16px
      &:hover
        color: $COLOR_BLUE
  .form
    padding: 30px
.bookFilter
  margin-bottom: 10px
  .categorySelector
    width: 110px
  .status
    width: 60px
  .divide
    padding-left: 30px
  .search
    padding-left: 40px
.table
  width: 100%
.pagination
  padding-top: 10px
  text-align: right
</style>
<script>
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
        if (!book.categories)  {
          return [];
        }
        return [allCategory].concat(book.categories);
      },
    })
  },
  methods: {
    ...mapActions([
      "bookList",
      "bookCacheRemove",
      "bookUpdate",
      "bookListCategory",
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
    edit(book) {
      this.currentUpdateBoook = Object.assign({}, book);
    },
    closeUpdate() {
      this.currentUpdateBoook = null;
    },
    async update() {
      const {
        currentUpdateBoook,
      } = this;
      const { id, statusDesc } = currentUpdateBoook;
      const found = find(this.currentBooks, item => item.id === id);
      if (!found) {
        return;
      }
      const updateData = {};
      const updateFields = [
        'brief',
        'name',
      ];
      updateFields.forEach((k) => {
        const v = currentUpdateBoook[k];
        if (found[k] !== v) {
          updateData[k] = v;
        }
      });
      if (found.statusDesc !== statusDesc) {
        updateData.status = this.bookStatusList.indexOf(statusDesc);
      }
      const currentCatgory = found.category.sort().join(',');
      const newCategory = currentUpdateBoook.category.sort().join(',');
      if (currentCatgory != newCategory) {
        updateData.category = newCategory;
      }
      const close = this.xLoading();
      try {
        await this.bookUpdate({
          id,
          update: updateData
        });
      } catch (err) {
        this.xError(err);
      } finally {
        this.currentUpdateBoook = null;
        close();
      }
    },
    async updateCover() {
      const {
        id
      } = this.currentUpdateBoook;
      const close = this.xLoading();
      try {
        await this.bookUpdateCover({
          id,
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
</script>
