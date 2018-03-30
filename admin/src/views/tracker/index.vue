<template lang="pug">
.trackerPage
  .form: el-form.clearfix.searchForm(
    v-model="form"
    label-width="90px"
  )
    el-form-item.pullLeft(
      label="Category"
      width="300"
    )
      el-select(
        v-model="form.category"
        placeholder="please select the category"
        clearable
      )
        el-option(
          v-for="item in categories"
          :key="item"
          :label="item"
          :value="item"
        )
    el-form-item.pullLeft(
      label="Account"
    )
      el-input.account(
        v-model="form.account"
        placeholder="search account"
        clearable
      )
    el-form-item.pullLeft(
      label="Begin"
      label-width="70px"
    )
      el-date-picker.datetimePicker(
        type="datetime"
        placeholder="Pick begin date"
        v-model="form.begin"
        value-format="timestamp"
      )
    el-form-item.pullLeft(
      label="End"
      label-width="70px"
    )
      el-date-picker.datetimePicker(
        type="datetime"
        placeholder="Pick end date"
        v-model="form.end"
        value-format="timestamp"
      )
    el-form-item.pullLeft(
      label-width="15px"
    )
      el-button(
        @click.native="search"
        type="primary"
      ) Search
  p.tac(
    v-if='trackers && trackers.length === 0'
  ) No result found. 
  el-table.mtop10(
    v-else-if='trackers'
    :data="trackers"
  )
    el-table-column(
      prop="category"
      label="Category"
      width="100"
    )
    el-table-column(
      prop="account"
      label="Account"
      width="100"
    )
    el-table-column(
      prop="ip"
      label="IP"
      width="150"
    )
    el-table-column(
      prop="params"
      label="Params"
    )
      template(
        slot-scope="scope"
      )
        pre {{scope.row.params}}
    el-table-column(
      prop="result"
      label="Result"
      width="80"
    )
    el-table-column(
      prop="use"
      label="Use"
      width="80"
    )
    el-table-column(
      prop="date"
      label="Date"
      width="160"
    )

</template>
<script src="./tracker.js"></script>
<style lang="sass" scoped>
@import "../../variables";
.trackerPage
  padding: 10px
.searchForm
  padding-top: 22px
  margin: auto
.account
  width: 140px
.datetimePicker
  width: 165px
</style>


