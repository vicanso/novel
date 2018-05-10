<template lang="pug">
  .bookList(v-if="mode === 0")
    p(v-if="books.length === 0") 很抱歉，没有符合条件的书籍
    div(
      v-else
    )
      el-form.form.searchForm.clearfix(
        v-model="searchForm"
        label-width="60px"
      )
        el-form-item.pullLeft(
          label="状态"
          width="200"
        )
          el-select(
            v-model="searchForm.status"
            placeholder="请选择书籍状态"
          )
            el-option(
              v-for="item in statusList"
              :key="item"
              :label="item"
              :value="item"
            )
        el-form-item.pullLeft.mleft10(
          label="关键字"
        )
          el-input(
            v-model="searchForm.keyword"
            placeholder="请输入关键字"
            clearable
          )
        el-form-item.pullLeft(
          label-width="30px"
        )
          el-button(
            @click.native="search"
            type="primary"
          ) 查询
        
      el-table(
        :data="books"
        stripe
      )
        el-table-column(
          prop="author"
          label="Author"
          width="100"
        )
        el-table-column(
          prop="name"
          label="Name"
          width="150"
        )
        el-table-column(
          label="End"
          width="80"
        )
          template(
            slot-scope="scope"
          )
            span(v-if="scope.row.end") YES
            span(v-else) NO
        el-table-column(
          label="Category"
          width="120"
          prop="category"
        )
        el-table-column(
          label="Latest Chapter"
          prop="latestChapter.title"
        )
        el-table-column(
          prop="date"
          label="Update At"
          width="160"
        )
        el-table-column(
          label="OP"
          width="220"
        )
          template(
            slot-scope="scope"
          )
            el-button(
              type="text"
              @click.native="edit(scope.row.no)"
            ) Edit
            el-button(
              type="text"
              @click.native="update(scope.row.no)"
            ) Update
            el-button(
              type="text"
              @click.native="updateCover(scope.row.no)"
            ) Update-Cover
      .clearfix.mtop10
        el-pagination.pullRight(
          background
          layout="prev, pager, next"
          :size="size"
          :total="count"
          @current-change="changePage"
        )
    el-button.addSource.mtop10(
      @click.native="addSource"
      type="primary"
    ) Add Source
  el-form.addForm.form(
    v-model="form"
    label-width="100px"
    v-else-if="mode === 1"
  )
    el-form-item(
      label="Author"
    )
      el-input(
        v-model="form.author"
        autofocus
        :disabled="true"
      )
    el-form-item(
      label="Name"
    )
      el-input(
        v-model="form.name"
        :disabled="true"
      )
    el-form-item(
      label="End"
    )
      el-switch(
        v-model="form.end"
      )
    el-form-item(
      label="Category"
    )
      el-select(
        v-model="form.category[0]"
      )
        el-option(
          v-for="item in categoryOptions"
          :key="item"
          :label="item"
          :value="item"
        )
    el-form-item(
      label="Brief"
    )
      el-input(
        type="textarea"
        :rows="8"
        v-model="form.brief"
      )
    el-form-item
      el-button(
        type="primary"
        @click.native="submit"
      ) Submit 
      el-button(
        @click.native="mode = 0"
      ) Back
  el-form.addForm.form(
    v-model="sourceForm"
    label-width="100px"
    v-else-if="mode === 2"
  )
    el-form-item(
      label="Source"
    )
      el-select(
        v-model="sourceForm.source"
      )
        el-option(
          v-for="item in sourceOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        )
    el-form-item(
      label="Author"
    )
      el-input(
        v-model="sourceForm.author"
      )
    el-form-item(
      label="Name"
    )
      el-input(
        v-model="sourceForm.name"
      )
    el-form-item(
      label="ID"
    )
      el-input(
        v-model="sourceForm.id"
      )
    el-form-item
      el-button(
        type="primary"
        @click.native="submitSource"
      ) Submit 
      el-button(
        @click.native="mode = 0"
      ) Back
</template>

<script src="./book.js"></script>
<style lang="sass" src="./book.sass" scoped></style>
