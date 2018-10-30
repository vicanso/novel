import localForage from "localforage";

const userSettingKey = "user-setting";

function getColors() {
  const shadow = "3px 6px 4px";
  const colors = {
    gray: {
      backgroundColor: "#d4d4d4",
      color: "#232323",
      boxShadow: `${shadow} rgba(125, 123, 116, 0.8)`
    },
    yellow: {
      backgroundColor: "#a89c84",
      color: "#4e3c26",
      boxShadow: `${shadow} rgba(125, 123, 116, 0.8)`
    },
    black: {
      backgroundColor: "#11100e",
      color: "#2d2c2b",
      boxShadow: `${shadow} rgba(30, 30, 30, 0.8)`
    },
    pinky: {
      backgroundColor: "#494446",
      color: "#2b0c12",
      boxShadow: `${shadow} rgba(60, 60, 60, 0.8)`
    }
  };
  return colors;
}

const defaultUserSetting = {
  fontSize: 22,
  colors: getColors(),
  theme: "yellow"
};

const chapterStore = localForage.createInstance({
  name: "chapters"
});

const bookReadInfoStore = localForage.createInstance({
  name: "book-read-info"
});

export async function clearChapterStoreExpired() {
  const keys = await chapterStore.keys();
  const now = Date.now();
  // 如果超过1个月则删除
  const ttl = 30 * 24 * 3600 * 1000;
  keys.forEach(async key => {
    const data = await chapterStore.getItem(key);
    const createdAt = data.createdAt || 0;
    if (now - createdAt > ttl) {
      await chapterStore.removeItem(key);
    }
  });
}

// ChapterCache 章节缓存
export class ChapterCache {
  constructor(id) {
    this.id = `${id}`;
  }
  async add(no, data) {
    const key = `${this.id}-${no}`;
    await chapterStore.setItem(key, {
      createdAt: Date.now(),
      data
    });
  }
  async get(no) {
    const key = `${this.id}-${no}`;
    const data = await chapterStore.getItem(key);
    if (!data) {
      return;
    }
    return data.data;
  }
}

// 获取已缓存的章节序号
export async function getStoreChapterIndexList(id) {
  const keys = await chapterStore.keys();
  const key = `${id}-`;
  const indexList = [];
  keys.forEach(v => {
    if (v.indexOf(key) !== -1) {
      const index = v.substring(key.length);
      indexList.push(Number.parseInt(index));
    }
  });
  return indexList.sort((a, b) => a - b);
}

// BookReadInfo 阅读信息
export class BookReadInfo {
  constructor(id) {
    this.id = `${id}`;
  }
  async update(no, page) {
    let data = await this.get();
    if (!data) {
      data = {
        createdAt: Date.now()
      };
    }
    Object.assign(data, {
      no,
      page,
      updatedAt: Date.now()
    });
    await bookReadInfoStore.setItem(this.id, data);
  }
  async get() {
    const data = await bookReadInfoStore.getItem(this.id);
    return data;
  }
}

export function saveUserSetting(data) {
  if (!localStorage) {
    return;
  }
  const currentSetting = getUserSetting();
  const value = JSON.stringify(Object.assign(currentSetting, data));
  localStorage.setItem(userSettingKey, value);
}

export function getUserSetting() {
  if (!localStorage) {
    return defaultUserSetting;
  }
  const data = localStorage.getItem(userSettingKey);
  let currentSetting = null;
  if (data) {
    currentSetting = JSON.parse(data);
  }
  return Object.assign({}, defaultUserSetting, currentSetting);
}
