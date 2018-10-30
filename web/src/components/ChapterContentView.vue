<template lang="pug">

mixin MainHeader
  mt-header.mainHeader(
    fixed
    v-if="isShowingSetting"
  )
    a.mainHeaderFunction(
      slot="left"
      @click="back"
    )
      i.mintui.mintui-back
mixin SettingFooter
  .settingFooter(
    v-if="isShowingSetting"
  )
    .fontSize
      span(
        @click="changeFontSize(-1)"
      ) T-
      span {{userSetting.fontSize}}
      span(
        @click="changeFontSize(1)"
      ) T+
    .theme
      div(
        v-for="theme, name in userSetting.colors"
        :key="name"
        @click="changeTheme(name)"
      )
        span(
          :style="{backgroundColor: theme.backgroundColor}"
        )
          i.iconfont.icon-selected(
            :class="{active: name === userSetting.theme}"
          )

.chapterContentWraper.fullHeight
  +MainHeader
  +SettingFooter
  .contentView.fullHeight(
    ref="contentView"
  ): .content(
    v-for="item in pages"
    :style="item.style"
  )
    h5 {{chapter.title}} 
    div(
      v-html="item.html"
    )
  .processBar.font12(
    v-if="chapterCount"
  ) {{(100 * chapterNo / chapterCount).toFixed(2) + "%"}}
</template>
<style lang="sass" scoped>
@import '@/styles/const.sass'
.content
  position: absolute
  left: 0
  top: 0
  right: 0
  bottom: 0
  padding-left: 15px
  h5
    line-height: 30px
    margin: 0
    padding: 0
.settingFooter
  position: fixed
  left: 0
  bottom: 0
  right: 0
  .fontSize
    line-height: 48px
    span
      text-align: center
      display: inline-block
      width: 30%
      &:nth-child(2)
        width: 40%
  .theme
    border-top: $GRAY_BORDER
    margin: 0 5px 
    padding-top: 5px
    div 
      display: inline-block
      width: 25%
    span
      display: block
      margin: 5px
      height: 40px
      line-height: 40px
      text-align: center
      i
        font-size: 32px
        visibility: hidden
        &.active
          visibility: visible
.mainHeader, .settingFooter
  background-color: rgba($COLOR_DARY_WHITE, 0.9)
  color: $COLOR_BLUE
  z-index: 99 !important
.processBar
  position: absolute
  bottom: 5px
  right: 15px
  z-index: 89
</style>

<script>
import { mapActions, mapState } from "vuex";
import Hammer from "hammerjs";
import { forEach } from "lodash";
import { getFontMetrics } from "@/helpers/util";

// 内容展示的padding
const padding = 15;
// 标题高度
const titleHeight = 30;
// 进度条高度
const progressBarHeight = 20;
// setting操作区
const settingFunctionsArea = {
  footerHeight: 110,
  top: 60,
  // 需要加载时计算
  bottom: 0
};

export default {
  props: {
    chapter: {
      type: Object,
      required: true
    },
    chapterNo: {
      type: Number,
      required: true
    },
    chapterCount: {
      type: Number
    },
    chapterPage: {
      type: Number
    }
  },
  data() {
    return {
      pages: null,
      currentPage: 0,
      maxPage: 0,
      isShowingSetting: false
    };
  },
  computed: {
    ...mapState({
      userSetting: ({ user }) => user.setting
    })
  },
  watch: {
    currentPage(v) {
      let index = 0;
      if (v <= 0) {
        // 切换至上一章
        index = -1;
      } else if (v > this.maxPage) {
        // 切换至下一章
        index = 1;
      }
      const chapterNo = this.chapterNo + index;
      if (chapterNo >= this.chapterCount) {
        this.xToast("已是最后一页");
        return;
      }
      if (chapterNo < 0) {
        this.xToast("已是第一页");
        return;
      }
      if (!index) {
        this.$emit("changePage", v);
        return;
      }
      this.pages = null;
      this.$emit("change", index);
    },
    async chapterNo(newValue, oldValue) {
      this.initPageContent();
      await this.$next();
      if (oldValue > newValue) {
        // 因为多了一页提示页，所以往前切换时，切换时直接用maxPage
        this.changePage(this.maxPage);
      } else {
        this.changePage(1);
      }
    }
  },
  methods: {
    ...mapActions(["userSaveSetting"]),
    initEvent() {
      const { $el, $refs } = this;
      const { maxWidth } = this.getOptions();
      const threshold = 5;
      const hammer = new Hammer($el, {
        direction: Hammer.DIRECTION_HORIZONTAL,
        threshold
      });
      const dom = $refs.contentView;
      // 移动的方向
      let direction = "";
      const directionLeft = "left";
      const directionRight = "right";
      const end = (item, currentPage, transX) => {
        this.currentPage = currentPage;
        const { style } = item;
        style.transition = "0.4s transform";
        style.transform = `translate3d(${transX}px, 0px, 0px)`;
      };
      hammer.on("pan panend panstart tap", e => {
        const { type, deltaX, center } = e;
        const { isShowingSetting } = this;
        if (type === "pen" && e.deltaX % 3 !== 0) {
          return;
        }
        if (type === "tap") {
          const { x, y } = center;
          // 如果已显示设置功能，tap则为取消
          if (isShowingSetting) {
            // 如果是顶部区域为操作区，忽略
            const { top, bottom } = settingFunctionsArea;
            if (y < top || y > bottom) {
              return;
            }

            this.isShowingSetting = false;
            return;
          }
          // 上一页
          if (x < 0.3 * maxWidth) {
            const v = this.currentPage - 1;
            end(dom.children[v], v, 0);
            return;
          }
          // 下一页
          if (x > 0.7 * maxWidth) {
            end(
              dom.children[this.currentPage],
              this.currentPage + 1,
              -maxWidth
            );
            return;
          }
          // 切换setting的显示
          this.isShowingSetting = !isShowingSetting;
          return;
        }
        if (type === "panstart") {
          if (deltaX > 0) {
            direction = directionRight;
          } else {
            direction = directionLeft;
          }
        }
        if (isShowingSetting) {
          return;
        }
        let index = 0;
        let offset = deltaX;
        if (direction === directionRight) {
          index -= 1;
          offset -= maxWidth;
        }
        const item = dom.children[this.currentPage + index];
        if (!item) {
          return;
        }
        switch (type) {
          case "panstart":
            item.style.transition = "";
            break;
          case "panend": {
            let currentPage = this.currentPage;
            // 如果移动小于30px，则认为无效操作
            const avalidMoving = Math.abs(deltaX) > 30;
            let transX = 0;
            if (direction === directionRight) {
              if (avalidMoving) {
                currentPage -= 1;
              } else {
                transX = -maxWidth;
              }
            } else if (avalidMoving) {
              transX = -maxWidth;
              currentPage += 1;
            }
            end(item, currentPage, transX);
            break;
          }
          default:
            item.style.transform = `translate3d(${offset +
              threshold}px, 0px, 0px)`;
            break;
        }
      });
      this.hammer = hammer;
    },
    getOptions() {
      const { $el, userSetting } = this;
      const { fontSize, colors, theme } = userSetting;
      const lineHeight = Math.ceil(fontSize * 1.5);
      return {
        maxWidth: $el.clientWidth + 10,
        lineHeight,
        fontSize,
        color: colors[theme]
      };
    },
    changePage(no) {
      const { $refs } = this;
      const { maxWidth } = this.getOptions();
      const dom = $refs.contentView;
      this.currentPage = no;
      forEach(dom.children, (item, index) => {
        const { style } = item;
        if (index < no) {
          style.transition = "";
          style.transform = `translate3d(${-maxWidth}px, 0px, 0px)`;
        } else {
          style.transition = "";
          style.transform = "translate3d(0px, 0px, 0px)";
        }
      });
    },
    initPageContent() {
      const { chapter, $el } = this;
      const { clientWidth, clientHeight } = $el;
      const { color, lineHeight, fontSize } = this.getOptions();
      const fontMetrics = getFontMetrics({
        width: clientWidth - 2 * padding,
        height: clientHeight - titleHeight - progressBarHeight,
        lineHeight: lineHeight,
        fontSize: fontSize,
        format: "html",
        color: color.color
      });
      const pages = fontMetrics.getFillTextList(chapter.content);
      this.maxPage = pages.length;
      pages.unshift({
        html: '<p class="tac">正在切换至上一章</p>'
      });
      pages.push({
        html: '<p class="tac">正在切换至下一章</p>'
      });
      const maxZIndex = pages.length;
      pages.forEach((item, i) => {
        item.style = Object.assign(
          {
            zIndex: maxZIndex - i
          },
          color
        );
      });
      this.pages = pages;
    },
    back() {
      this.$emit("back");
    },
    changeFontSize(v) {
      const { userSetting } = this;
      userSetting.fontSize += v;
      this.userSaveSetting(userSetting);
      this.initPageContent();
    },
    changeTheme(name) {
      const { userSetting } = this;
      if (name === userSetting.theme) {
        return;
      }
      userSetting.theme = name;
      this.userSaveSetting(userSetting);
      this.initPageContent();
    }
  },
  async mounted() {
    this.initPageContent();
    await this.$next();
    let chapterPage = 1;
    if (this.chapterPage) {
      chapterPage = this.chapterPage;
    }
    this.changePage(chapterPage);
    this.initEvent();
    settingFunctionsArea.bottom =
      this.$el.clientHeight - settingFunctionsArea.footerHeight;
  },
  beforeDestroy() {
    this.hammer.destroy();
  }
};
</script>
