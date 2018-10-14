<template lang="pug">
.home
  h5 the right one is the best
</template>

<style lang="sass" scoped>
@import "@/styles/const.sass"
.home
  padding: 30px 80px
  color: $COLOR_BLACK
  ul, p
    font-size: 14px
h5, h6
  font-family: Oleo Script,cursive
  text-align: center
h5
  font-size: 32px
  text-transform: uppercase
h6
  font-size: 18px
</style>


<script>
export default {
  name: "home",
  components: {}
};
</script>
