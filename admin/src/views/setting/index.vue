<template lang="pug">
  .settingPage 
    h3.tac System Setting 
    div(
      v-if="mode === 0"
    )
      p.tac(v-if="!settings || settings.length === 0") There isn't any system setting
      el-table.table(
        v-else
        :data="settings"
      )
        el-table-column(
          prop="name"
          label="Name"
          width="150"
        )
        el-table-column(
          label="Data"
        )
          template(
            slot-scope="scope"
          )
            pre {{scope.row.data}}
        el-table-column(
          label="Status"
          width="100"
        )
          template(
            slot-scope="scope"
          )
            span(v-if="scope.row.disabled") Disabled
            span(v-else) Enabled
        el-table-column(
          prop="description"
          label="Description"
          width="200"
        )
        el-table-column(
          prop="date"
          label="Date"
          width="140"
        )
        el-table-column(
          prop="creator"
          label="Creator"
          width="100"
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
      ) Add Setting
    el-form.form(
      v-model="form"
      label-width="100px"
      v-if="mode === 1"
    )
      el-form-item(
        label="Name"
      )
        el-input(
          v-model="form.name"
          autofocus
        )
      el-form-item(
        label="Disabled"
      )
        el-switch(
          v-model="form.disabled"
        )
      el-form-item(
        label="Data"
      )
        el-input(
          type="textarea"
          v-model="form.data"
          :rows="8"
        )
      el-form-item(
        label="Description"
      )
        el-input(
          type="textarea"
          v-model="form.description"
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
<script src="./setting.js"></script>
<style lang="sass" scoped>
.settingPage
  padding: 10px
.form
  width: 600px
  padding: 40px 20px 20px
  margin: auto
.addSetting, .table
  width: 100%
</style>

