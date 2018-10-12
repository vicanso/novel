<template lang="pug">
.bookList(
  ref="bookList"
)
  .tableWrapper(
    v-if="!loading"
  )
    el-table.table(
      :data="currentBooks"
      :height="tableHeight"
      border
      stripe
    )
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
.table
  width: 100%
.pagination
  padding: 15px 0
  text-align: right
</style>
<script>
export default {
  name: 'book-list',
  data() {
    return {
      currentBooks: null,
      loading: false,
      tableHeight: 0,
    };
  },
  mounted() {
    this.tableHeight = this.$refs.bookList.clientHeight - 80;
  },
}
</script>
