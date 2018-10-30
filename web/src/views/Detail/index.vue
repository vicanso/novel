<template lang="pug">

mixin MainContent
  .content(
    v-if="detail"
  )
    .mainInfos.clearfix
      .cover.pullLeft(
        v-if="cover"
      )
        ImageView(
          :src="cover"
        )
      .infos.font14
        span.font12(
          v-for="item in detail.category"
        ) {{item}}
        span.font12(
          v-if="detail.category && detail.category.length"
        ) | 
        span.font12 {{detail.author}}
        .wordCount {{wordCountDesc}}
        .latestChapter.font12(
          v-if="latestChapter"
        ) 最新章节：
          span {{latestChapter.title}}
    .brief.font14 {{detail.brief}}
    .chapters(
      v-if="chapterCount"
      @click="showChapters()"
    ) 查看目录
      i.pullRight.mintui.mintui-back.rotate180.font14
      span.pullRight.font12.updatedAt 更新于
        span {{latestChapter.updatedAt.substring(5, 16)}} 
      span.mleft5.font12 连载至{{chapterCount}}章

mixin Recommends
  .recommends(
    v-if="recommendBooks && recommendBooks.length"
  )
    h3 同类热门
    BookView.book(
      v-for="item in recommendBooks"
      :key="item.id"
      :id="item.id"
      :name="item.name"
      :author="item.author"
      :brief="item.brief"
      :cover="item.cover"
      :wordCount="item.wordCount"
    )

mixin MainView
  mt-header.mainHeader(
    :title="(detail && detail.name) || '...'"
    fixed
  )
    a.mainHeaderFunction(
      slot="left"
      @click="back('')"
    )
      i.mintui.mintui-back
    a.refresh.mainHeaderFunction(
      slot="right"
      @click="refresh()"
    )
      i.iconfont.icon-refresh
  .contentWrapper.fullHeightScroll
    +MainContent
    +Recommends
  .functions(
    v-if="detail"
  )
    a(
      href="javascript:;"
      @click="download()"
    )
      i.iconfont.icon-icondownload
      span 下载
    a.read(
      href="javascript:;"
      @click="startReading()"
    )
      i.iconfont.icon-office
      span(
        v-if="!read"
      ) 免费阅读
      span(
        v-else
      ) 继续阅读
    a(
      href="javascript:;"
      @click="addToShelf()"
    )
      i.iconfont.icon-pin
      span 加入书架

mixin ChaptersView
  .fullHeight.chaptersView
    mt-header.mainHeader(
      :title="detail && detail.name"
      fixed
    )
      a.mainHeaderFunction(
        slot="left"
        @click="back('main')"
      )
        i.mintui.mintui-back
    .chaptersWrapper.fullHeightScroll
      ul.chapterSection
        li(
          v-for="item, index in chapterSections"
          :key="item.start"
        ): a(
          :class="currentChapterSection == index ? 'active' : ''"
          href="javascript:;"
          @click="changeChapterSection(index)"
        ) {{item.start + "-" + item.end}}
        li.sortBy
          a(
            @click="toggleChapterOrder"
            href="javascript:;"
          ) 倒序
      ul.chapters
        li(
          v-for="item, index in currentChapters"
          :key="item.index"
        ): a(
          href="javascript:;"
          @click="startToReadChapter(item.no)"
        ) 
          span.pullRight.font12(
            v-if="isStored(item.no)"
          ) 已下载
          | {{item.title}}

mixin ReadChapterView
  .fullHeight.readChapterView(
    :style="readChapterViewStyle"
  )
    ChapterContentView(
      v-if="currentChapter"
      :chapter="currentChapter"
      :chapterNo="currentChapterNo"
      :chapterPage="currentChapterPage"
      :chapterCount="chapterCount"
      @change="changeChapter"
      @changePage="changeChapterPage"
      @back="backFromRead"
    )

.detailWrapper.fullHeight
  .fullHeight(
    v-show="view === 'chapters'"
  )
    +ChaptersView
  .fullHeight(
    v-if="view === 'readChapter'"
  )
    +ReadChapterView
  .fullHeight.mainView(
    v-show="view === 'main'"
  )
    +MainView
</template>

<script src="./detail.js"></script>
<style lang="sass" src="./detail.sass" scoped></style>

