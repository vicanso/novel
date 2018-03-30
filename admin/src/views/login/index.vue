<template lang="pug">
  .loginPage
    el-form.form(
      ref="form"
      v-model="form"
      label-width="100px"
    )
      h3.tac Login
      el-form-item(
        label="Account"
      )
        el-input(
          v-model="form.account"
          @keyup.enter.native="submit"
          autofocus
        )
      el-form-item(
        label="Password"
      )
        el-input(
          v-model="form.password"
          @keyup.enter.native="submit"
          type="password"
          auto-complete="off"
        )
      el-form-item(style="margin-bottom:0")
        el-button(
          type='primary'
          @click="submit"
        ) Submit
</template>
<style lang="sass" scoped>
.loginPage
  padding: 10px
  width: 600px
  margin: auto
.form
  margin-top: 60px
  padding: 10px 30px 30px
</style>
<script src="./login.js"></script>

