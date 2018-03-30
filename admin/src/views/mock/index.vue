<template lang="pug">
  .mockPage
    h3.tac Mock Setting
    div(
      v-if="mode === 0"
    )
      p.tac(v-if="!mocks || mocks.length === 0") There isn't any mock setting
      el-table.table(
        v-else
        :data="mocks"
      )
        el-table-column(
          prop="account"
          label="Account"
          width="100"
        ) 
        el-table-column(
          prop="url"
          label="Url"
          width="150"
        )
        el-table-column(
          prop="status"
          label="Status"
          width="70"
        )
        el-table-column(
          label="Response"
        )
          template(
            slot-scope="scope"
          )
            pre {{scope.row.response}}
        el-table-column(
          label="Disabled"
        )
          template(
            slot-scope="scope"
          )
            span(v-if="scope.row.disabled") Disabled
            span(v-else) Enabled
        el-table-column(
          label="Description"
          prop="description"
        )
        el-table-column(
          label="OP"
          width="80"
        )
          template(
            slot-scope="scope"
          )
            el-button(
              type="text"
              @click.native="edit(scope.row.id)"
            ) Edit
      el-button.addSetting.mtop10(
        type="primary"
        @click.native="add"
      ) Add Mock
    el-form.form(
      v-model="form"
      label-width="100px"
      v-if="mode === 1"
    )
      el-form-item(
        label="Account"
      )
        el-input(
          v-model="form.account"
          autofocus
        )
      el-form-item(
        label="Url"
      )
        el-input(
          v-model="form.url"
        )
      el-form-item(
        label="Status"
      )
        el-input(
          type="number"
          v-model="form.status"
        )
      el-form-item(
        label="Response"
      )
        el-input(
          v-model="form.response"
          type="textarea"
          :rows="8"
        )
      el-form-item(
        label="Disabled"
      )
        el-switch(
          v-model="form.disabled"
        )
      el-form-item(
        label="Description"
      )
        el-input(
          v-model="form.description"
          type="textarea"
          :rows="4"
        )
      el-form-item
        el-button(
          type="primary"
          @click.native="submit"
        ) Submit 
        el-button(
          @click.native="mode = 0"
        ) Back
</template>
<script src="./mock.js"></script>
<style lang="sass" scoped>
.mockPage
  padding: 10px
.form
  width: 600px
  margin: auto
  padding: 40px 20px 20px
  margin: auto
.addSetting, .table
  width: 100%
</style>


