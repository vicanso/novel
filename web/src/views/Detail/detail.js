import { mapActions, mapState } from "vuex";
import ImageView from "@/components/ImageView";
import BookView from "@/components/BookView";
import ChapterContentView from "@/components/ChapterContentView";
import {
  getCover,
  formatDate,
  waitfor,
} from "@/helpers/util";
import {
  routeDetail,
} from "@/routes";

const sectionChapterCount = 100
const readChapterView = "readChapter";
const mainView = "main";
const chaptersView = "chapters";


export default {
  name: "detail",
  data() {
    return {
      id: "",
      view: mainView,
      prevView: "",
      cover: "",
      wordCountDesc: "",
      detail: null,
      // 章节总数
      chapterCount: 0,
      // 最新章节
      latestChapter: null,
      // 相关推荐
      recommendBooks: null,
      // 章节分组
      chapterSections: null,
      // 当时章节分组
      currentChapterSection: -1,
      currentChapters: null,
      // 当前章节标题、内容
      currentChapter: null,
      // 当前章节序号
      currentChapterNo: 0,
      // 当前章节页
      currentChapterPage: 0,
      // 章节排序（升序）
      chapterOrderAsc: true,
      // 是否已阅读过本书
      read: false,
      // 已缓存的章节序号
      storeChapterInfos: null,
    };
  },
  components: {
    ImageView,
    BookView,
    ChapterContentView,
  },
  computed: {
    ...mapState({
      readChapterViewStyle: ({ user }) => {
        if (!user.setting) {
          return null;
        }
        const {
          colors,
          theme,
        } = user.setting;
        if (!colors || !theme) {
          return null;
        }
        const v = colors[theme];
        if (!v) {
          return null;
        }
        return {
          "backgroundColor": v.backgroundColor,
        };
      },
    }),
  },
  watch: {
    detail(v) {
      if (!v) {
        return;
      }
      const {
        cover,
        wordCount,
      } = v;
      this.cover = getCover(cover, 100);
      const base = 10 * 1000;
      let wordCountDesc = "";
      if (wordCount >= base) {
        wordCountDesc = `${Math.floor(wordCount / base)}万字`;
      } else {
        wordCountDesc = `${wordCount}字`;
      }
      this.wordCountDesc = wordCountDesc;
    },
    view(v) {
      if (v !== readChapterView) {
        // 清除数据
        this.currentChapter = null;
      }
      // 如果是章节内容阅读，则不记录
      if (v === readChapterView) {
        return;
      }
      this.prevView = v;
    },
  },
  methods: {
    ...mapActions([
      "bookGetDetail",
      "bookGetRecommend",
      "bookGetChapters",
      "bookGetChapterContent",
      "bookGetReadInfo",
      "bookUpdateReadInfo",
      "bookGetStoreChapterIndexes",
      "bookDownload",
    ]),
    async load(id) {
      const close = this.xLoading();
      try {
        const res = await this.bookGetDetail({
          id,
        });
        const {
          book,
          chapterCount,
          latestChapter,
        } = res.data;
        const {
          cover,
          wordCount,
        } = book;
        this.cover = getCover(cover, 100);
        const base = 10 * 1000;
        let wordCountDesc = "";
        if (wordCount >= base) {
          wordCountDesc = `${Math.floor(wordCount / base)}万字`;
        } else {
          wordCountDesc = `${wordCount}字`;
        }
        this.wordCountDesc = wordCountDesc;
        this.detail = book;
        this.chapterCount = chapterCount;
        if (latestChapter) {
          latestChapter.updatedAt = formatDate(latestChapter.updatedAt);
        }
        const chapterSectionCount = Math.ceil(chapterCount / sectionChapterCount);
        const chapterSections = [];
        for (let i = 0; i < chapterSectionCount; i++) {
          const start = i * sectionChapterCount + 1;
          let end = (i + 1) * sectionChapterCount;
          if (end > chapterCount) {
            end = chapterCount;
          }
          chapterSections.push({
            start,
            end,
          });
        }
        this.chapterSections = chapterSections;
        this.latestChapter = latestChapter;
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    },
    back(view) {
      if (!view) {
        this.$router.back();
        return
      }
      this.view = view;
    },
    async loadRecommend(id) {
      try {
        const res = await 
        this.bookGetRecommend({
          id,
          limit: 3,
          field: "id,name,author,brief,cover,wordCount",
          order: "-updatedAt",
        });
        this.recommendBooks = res.data.books;
      } catch (err) {
        this.xError(err);
      }
    },
    async changeChapterSection(index) {
      const {
        currentChapterSection,
        chapterSections,
        id,
      } = this;
      if (currentChapterSection == index) {
        return
      }
      this.currentChapterSection = index;
      const data = chapterSections[index];
      const close = this.xLoading();
      try {
        const offset = data.start - 1;
        const res = await this.bookGetChapters({
          id,
          limit: sectionChapterCount,
          offset,
          order: "index",
          field: "title,updatedAt,index",
        });
        const {
          chapters,
        } = res.data;
        chapters.forEach((v, i) => {
          v.no = offset + i;
        })
        if (this.chapterOrderAsc) {
          this.currentChapters = chapters;
        } else {
          this.currentChapters = chapters.reverse();
        }
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    },
    async showChapters() {
      this.view = chaptersView;
      try {
        const indexList = await this.bookGetStoreChapterIndexes({
          id: this.id,
        });
        const storeChapterInfos = {};
        indexList.forEach((v) => {
          storeChapterInfos[v] = true;
        });
        this.storeChapterInfos = storeChapterInfos
      } catch (err) {
        this.xError(err);
      }
      if (this.currentChapterSection === -1) {
        this.changeChapterSection(0);
      }
    },
    isStored(no) {
      const {
        storeChapterInfos,
      } = this;
      return storeChapterInfos[no];
    },
    reset() {
      this.detail = null;
      this.currentChapterSection = -1;
      this.currentChapters = null;
      this.recommendBooks = null;
    },
    async init(id) {
      this.reset();
      this.id = id;
      const data = await this.bookGetReadInfo({
        id,
      });
      if (data) {
        this.currentChapterNo = data.no || 0;
        this.currentChapterPage = data.page || 0;
        this.read = true;
      }
      await this.load(id);
      await this.loadRecommend(id);
    },
    async showChapterContent(no, page) {
      const {
        id,
        view,
      } = this;
      const close = this.xLoading();
      if (view !== readChapterView) {
        this.view = readChapterView;
      }
      try {
        const done = waitfor(300);
        const data = await this.bookGetChapterContent({
          id,
          no,
        });
        await this.bookUpdateReadInfo({
          id,
          no,
          page,
        });
        await done();
        this.currentChapter = data;
        this.currentChapterNo = no;
        this.currentChapterPage = page;
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    },
    async startReading() {
      const {
        id,
        currentChapterNo,
        currentChapterPage,
      } = this;
      const data = await this.bookGetReadInfo({
        id,
      });
      let no = currentChapterNo;
      let page = currentChapterPage;
      if (data) {
        no = data.no;
        page = data.page;
      }
      this.showChapterContent(no, page);
    },
    startToReadChapter(no) {
      this.showChapterContent(no, 0)
    },
    changeChapter(index) {
      this.showChapterContent(this.currentChapterNo + index, 0);
    },
    changeChapterPage(page) {
      const {
        id,
        currentChapterNo,
      } = this;
      this.bookUpdateReadInfo({
        id,
        no: currentChapterNo,
        page,
      });
    },
    backFromRead() {
      this.view = this.prevView || mainView;
    },
    toggleChapterOrder() {
      const {
        chapterOrderAsc,
        chapterSections,
        currentChapterSection,
      } = this;
      this.currentChapters = null;
      this.chapterOrderAsc = !chapterOrderAsc;
      this.chapterSections = chapterSections.reverse();
      // 重置
      this.currentChapterSection = -1;
      this.changeChapterSection(currentChapterSection);
    },
    refresh() {
      const {id} = this.$route.params;
      this.init(id);
    },
    async download() {
      const close = this.xLoading({
        timeout: 300 * 1000,
      });
      try {
        await this.bookDownload({
          id: this.id,
          max: this.chapterCount,
        });
        this.xToast("已全部下载完成");
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    },
    addToShelf() {
      this.xToast("敬请期待");
    }
  },
  beforeMount() {
    this.refresh();
  },
  beforeRouteUpdate(to, from, next) {
    const {name, params} = to;
    if (name === routeDetail) {
      const {id} = params;
      this.init(id);
    }
    next();
  },
}