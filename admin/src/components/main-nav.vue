<template lang="pug">
  .mainNav
    el-menu(
      :default-active="active"
    )
      el-menu-item(
        v-for="(item, index) in items",
        :key="index",
        :index="'' + index"
        @click.native="go(item)"
      )
        span {{item.name}}
</template>

<style lang="sass" scoped>
@import "../variables";
.mainNav
  position: fixed    
  left: 0
  top: $MAIN_HEADER_HEIGHT
  bottom: 0
  width: $MAIN_NAV_WIDTH
  background-color: $COLOR_BLACK
  z-index: 9
</style>
<script>
export default {
  data() {
    return {
      active: '0',
      items: [
        {
          name: 'System Setting',
          route: 'setting',
        },
        {
          name: 'Mock',
          route: 'mock',
        },
        {
          name: 'Tracker',
          route: 'tracker',
        },
      ],
    };
  },
  methods: {
    go(item) {
      this.$router.push({
        name: item.route,
      });
    },
  },
};
</script>
