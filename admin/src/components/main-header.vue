<template lang="pug">
  .mainHeader
    .pullRight(style="margin-right:15px")
      span(v-if='!userInfo') Loading...
      div(
        v-else-if='userInfo.anonymous'
      )
        el-button(
          type="text"
          @click="$router.push({name: 'register'})"
        ) register
        el-button(
          type="text"
          @click="$router.push({name: 'login'})"
        ) login
      div(v-else)
        span.mright10 {{userInfo.account}}
        el-button(
          type="text"
          @click="userLogout"
        ) logout
    .logo Novel
</template>
<style lang="sass" scoped>
@import "../variables";
.mainHeader
  position: fixed
  left: 0
  top: 0
  right: 0
  height: $MAIN_HEADER_HEIGHT
  background-color: $COLOR_WHITE
  border-bottom: $GRAY_BORDER
  z-index: 9
  line-height: $MAIN_HEADER_HEIGHT
  .logo
    width: $MAIN_NAV_WIDTH
    background-color: $COLOR_BLACK
    color: $COLOR_WHITE
    text-align: center
    font-size: 18px
</style>

<script>
import {mapState, mapActions} from 'vuex';
export default {
  methods: {
    ...mapActions(['userLogout']),
  },
  computed: {
    ...mapState({
      userInfo: ({user}) => user.info,
    }),
  },
};
</script>
