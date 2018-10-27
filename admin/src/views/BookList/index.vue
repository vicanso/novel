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
        label="最新章节"
        v-if="currentUpdateBoook.latestChapter"
      )
        el-input(
          v-model="currentUpdateBoook.latestChapter.title"
          :disabled="true"
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
          @click.native="updateCover"
        ) 更新
        el-input(
          v-model="currentUpdateBoook.sourceCover"
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
  z-index: 1
.bookUpdate
  position: fixed 
  top: 0
  left: 0
  right: 0
  bottom: 0
  background-color: $COLOR_WHITE
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
<script src="./bookList.js">
</script>
