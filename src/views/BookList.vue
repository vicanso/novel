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
    h5 编辑内容
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
          :disabled="true"
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
        label="原始封面"
      )
        img(
          :src="currentUpdateBoook.sourceCover"
          height="90px"
        ) 
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
        ) 保存


.bookList(
  ref="bookList"
)
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
  padding: 15px
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
  .form
    padding: 30px
.table
  width: 100%
.pagination
  padding: 15px 0
  text-align: right
</style>
<script>
import { mapActions, mapState } from "vuex";
import { getListBookPageSize, saveListBookPageSize } from "@/helpers/storage";
export default {
  name: "book-list",
  data() {
    return {
      currentBooks: null,
      loading: true,
      tableHeight: 0,
      field: ["name", "author", "brief", "status", "updatedAt", "sourceCover"].join(","),
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
      bookStatusList: ( { book }) => book.statusList,
    })
  },
  methods: {
    ...mapActions(["bookList", "bookCacheRemove"]),
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
    async fetch() {
      const { field, order, offset, limit } = this;
      const close = this.xLoading();
      this.loading = true;
      try {
        await this.bookList({
          field,
          order,
          offset,
          limit
        });
        this.currentBooks = this.books.slice(offset, offset + limit);
      } catch (err) {
        this.xError(err);
      } finally {
        this.loading = false;
        close();
      }
    },
    async edit(book) {
      this.currentUpdateBoook = book;
    },
  },
  mounted() {
    this.tableHeight = this.$refs.bookList.clientHeight - 80;
  },
  beforeMount() {
    this.fetch();
  }
};
</script>
