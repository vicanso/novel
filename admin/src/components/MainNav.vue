<template lang="pug">
nav.nav
  ul(
    v-for="item in categories"
  )
    li: a(
      :class="active == item.name ? 'active': ''"
      href="javascript:;"
      @click="go(item)"
    ) {{item.name}}
</template>

<style lang="sass" scoped>
@import '@/styles/const.sass'
.nav
  position: fixed    
  left: 0
  top: $MAIN_HEADER_HEIGHT
  bottom: 0
  width: $MAIN_NAV_WIDTH
  background-color: $COLOR_BLACK
  z-index: 9
ul
  $height: 50px
  margin: 0
  padding: 0
  list-style: none
  line-height: $height
  a
    display: block
    color: $COLOR_WHITE
    text-decoration: none
    padding-left: 15px
    border-left: 3px solid $COLOR_BLACK
    &:hover, &.active
      color: $COLOR_BLUE
      background-color: $COLOR_WHITE
      border-left-color: $COLOR_BLUE
</style>

<script>
import { routeBookList } from "@/routes";
export default {
  name: "main-nav",
  data() {
    return {
      active: this.$route.params.category || "",
      categories: [
        {
          name: "books",
          path: routeBookList
        }
      ]
    };
  },
  methods: {
    go({ name, path }) {
      this.active = name;
      this.$router.push({
        name: path
      });
    }
  }
};
</script>
