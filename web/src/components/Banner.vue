<template lang="pug">
.banner.clearfix(
  ref="banner"
  :style="style"
)
  ImageView.pullLeft(
    v-for="item in images"
    :key="item"
    :src="item"
  )
  .dots
    span(
      v-for="item, index in images"
      :class="active == index ? 'active' : ''"
    )
</template>

<style lang="sass" scoped>
@import '@/styles/const.sass'
.banner
  height: 100%
  position: relative
.dots
  position: absolute
  left: 10px
  bottom: 5px
  span
    display: inline-block
    background-color: $COLOR_GRAY
    width: 5px
    height: 5px
    border-radius: 10px
    margin-right: 3px
    &.active
      background-color: $COLOR_WHITE
</style>

<script>
import ImageView from "@/components/ImageView";
export default {
  components: {
    ImageView
  },
  name: "banner",
  props: {
    ids: Array
  },
  data() {
    return {
      images: [],
      style: {},
      active: 0
    };
  },
  mounted() {
    const baseUrl =
      "http://tiny.aslant.site/api/images/v1/center/:id-90-:w-:h.jpeg";
    const dom = this.$refs.banner;
    const w = dom.clientWidth;
    const h = dom.clientHeight;
    this.ids.forEach(id => {
      const url = baseUrl
        .replace(":id", id)
        .replace(":w", w)
        .replace(":h", h);
      this.images.push(url);
    });
    this.style.width = `${this.images.length}00%`;
  }
};
</script>
