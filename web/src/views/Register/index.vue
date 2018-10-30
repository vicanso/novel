<template lang="pug">
.registerWrapper.fullHeight
  .tips 每个人都可以选择自己的生活，不骄不躁，做最平凡的自己！
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
  mt-field(
    label="确认密码"
    type="password"
    placeholder="请再次输入您的密码"
    v-model.trim="pwd"
  )
  mt-field(
    label="邮箱"
    type="email"
    placeholder="请输入您的邮箱"
    v-model.trim="email"
  )
  mt-button.login(
    @click.native="register"
    type="primary"
  ) 立即注册
</template>

<script src="./register.js"></script>
<style lang="sass" scoped src="./register.sass"></style>
