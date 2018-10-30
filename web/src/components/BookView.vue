<template lang="pug">
.bookView(
  @click="showDetail"
)
  .cover.pullLeft
    ImageView(
      v-if="coverUrl"
      :src="coverUrl"
    )
  .content
    h5 {{name}}
    p {{brief}}
    .otherInfos
      .tags.pullRight.font12
        span(
          v-if="categoryDesc"
        ) {{categoryDesc}}
        span(
          v-if="wordCountDesc"
        ) {{wordCountDesc}}
      .author
        i.iconfont.icon-account
        | {{author}}
</template>

<style lang="sass" scoped>
@import '@/styles/const.sass'
$coverWidth: 80px
.bookView
  padding: 10px
.cover
  width: $coverWidth
  height: 100%
  overflow: hidden
.content
  margin-left: $coverWidth + 10 
  $lineHeight: 20px
  line-height: $lineHeight 
  h5
    margin: 0
    padding: 0
    line-height: 2em
    font-size: 16px
  p
    font-size: 13px
    margin: 0 0 10px 0
    height: 2 * $lineHeight 
    overflow: hidden
    color: $COLOR_DARK_GRAY
    text-overflow: ellipsis
    word-break: break-all
  .otherInfos
    height: $lineHeight
    overflow: hidden
  .author
    color: $COLOR_DARK_GRAY
    font-size: 13px
    i
      margin-right: 5px
      font-weight: 900
      font-size: 10px
      color: $COLOR_WHITE
      background-color: $COLOR_DARK_BLUE
      border-radius: 10px
      padding: 2px
  .tags
    margin-top: 1px
    span
      border-radius: 3px
      margin: 0
      margin-left: 5px
      padding: 1px 4px
      &:nth-child(2n)
        color: $COLOR_DARK_GRAY
        border: 1px solid $COLOR_DARK_GRAY
      &:nth-child(2n + 1)
        color: $COLOR_DARK_BLUE
        border: 1px solid $COLOR_DARK_BLUE
</style>


<script>
import ImageView from "@/components/ImageView";
import { routeDetail } from "@/routes";
import { getCover } from "@/helpers/util";

const ignoreCategory = ["今日必读"];

const coverHeight = 98;
export default {
  components: {
    ImageView
  },
  props: {
    id: {
      type: Number,
      required: true
    },
    name: {
      type: String,
      required: true
    },
    author: {
      type: String,
      required: true
    },
    brief: {
      type: String,
      required: true
    },
    wordCount: {
      type: Number
    },
    cover: {
      type: String
    },
    category: {
      type: Array
    }
  },
  data() {
    let wordCountDesc = "";
    const { wordCount, category } = this;
    if (wordCount) {
      const tenThousand = 10000;
      if (wordCount >= tenThousand) {
        wordCountDesc = `${Math.floor(wordCount / tenThousand)}万字`;
      } else {
        wordCountDesc = `${wordCount}字`;
      }
    }
    let categoryDesc = "";
    if (category) {
      category.forEach(v => {
        if (!categoryDesc && ignoreCategory.indexOf(v) === -1) {
          categoryDesc = v;
        }
      });
    }

    return {
      coverUrl: "",
      wordCountDesc,
      categoryDesc
    };
  },
  mounted() {
    const { cover } = this;
    if (!cover) {
      return;
    }
    this.coverUrl = getCover(cover, coverHeight);
  },
  methods: {
    showDetail() {
      this.$router.push({
        name: routeDetail,
        params: {
          id: this.id
        }
      });
    }
  }
};
</script>
