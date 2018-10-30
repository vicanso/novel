<template lang="pug">
.imageWrapper(
  v-intersecting-show="load"
)
  img(
    v-if="loaded"
    :src="imageSrc"
  )
</template>

<style lang="sass" scoped>
.imageWrapper
  height: 100%
  img
    display: block
    margin: auto
</style>


<script>
import { supportWebp } from "@/helpers/util";

export default {
  props: {
    // 图片地址，如果指定了图片地址，直接加载
    src: {
      type: String
    }
  },
  data() {
    return {
      loaded: false,
      imageSrc: ""
    };
  },
  methods: {
    load() {
      const { src } = this;
      let imageSrc = src;
      if (supportWebp()) {
        imageSrc = src.replace(/\.png|\.jpeg/, ".webp");
      }
      const img = new Image();
      img.onload = () => {
        this.loaded = true;
      };
      this.imageSrc = imageSrc;
      img.src = imageSrc;
    }
  }
};
</script>
