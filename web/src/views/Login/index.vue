<template lang="pug">
.loginWrapper.fullHeight
  .tips 每个人都有自己的成功之道，不妒不忌，做最真实的自己！
  mt-field(
    label="用户名"
    placeholder="请输入您的用户名"
    v-model.trim="account"
  )
  mt-field(
    label="密码"
    type="password"
    placeholder="请输入您的密码"
    v-model.trim="password"
  )
  mt-button.login(
    @click.native="login"
    type="primary"
  ) 立即登录
</template>

<script src="./login.js"></script>
<style lang="sass" scoped src="./login.sass"></style>


