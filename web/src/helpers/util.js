import FontMetrics from "web-font-metrics";
import { isError, map } from "lodash-es";
import { env, coverUrlPrefix } from "@/config";

let currentFontMetrics = null;

export function log(...args) {
  // eslint-disable-next-line
  console.info(...args);
}

export function isDevelopment() {
  return env == "development";
}

export function debug(...args) {
  if (!isDevelopment()) {
    return;
  }
  // eslint-disable-next-line
  console.debug(...args);
}

// 获取出错信息
export function getErrorMessage(err) {
  let message = err;
  if (err && err.response) {
    const { data, headers } = err.response;
    if (data.message) {
      message = data.message;
    }
    const id = headers["x-response-id"];
    if (data.code) {
      // eslint-disable-next-line
      const code = data.code.replace(`${app}-`, "");
      message = `${data.message}(${code}) [${id}]`;
    }
  }
  if (isError(message)) {
    message = message.message;
  }
  if (err.code === "ECONNABORTED") {
    message = "请求超时，请重新再试";
  }
  return message;
}

// formatDate
export function formatDate(str) {
  if (!str) {
    return "";
  }
  const date = new Date(str);
  const fill = v => {
    if (v >= 10) {
      return `${v}`;
    }
    return `0${v}`;
  };
  const month = fill(date.getMonth() + 1);
  const day = fill(date.getDate());
  const hours = fill(date.getHours());
  const mintues = fill(date.getMinutes());
  const seconds = fill(date.getSeconds());
  return `${date.getFullYear()}-${month}-${day} ${hours}:${mintues}:${seconds}`;
}

// copy copy the value
export function copy(value) {
  // 来源自：https://juejin.im/post/5a94f8eff265da4e9b593c29
  const input = document.createElement("input");
  input.setAttribute("readonly", "readonly");
  input.setAttribute("value", value);
  document.body.appendChild(input);
  input.select();
  input.setSelectionRange(0, 9999);
  document.execCommand("copy");
  document.body.removeChild(input);
}

let isSupportWebp = false;
(function() {
  const images = {
    basic:
      "data:image/webp;base64,UklGRjIAAABXRUJQVlA4ICYAAACyAgCdASoCAAEALmk0mk0iIiIiIgBoSygABc6zbAAA/v56QAAAAA==",
    lossless:
      "data:image/webp;base64,UklGRh4AAABXRUJQVlA4TBEAAAAvAQAAAAfQ//73v/+BiOh/AAA="
  };
  const check = data =>
    new Promise((resolve, reject) => {
      const img = new Image();
      img.onload = resolve;
      img.onerror = reject;
      img.src = data;
    });
  Promise.all(map(images, check))
    .then(() => true)
    .catch(() => false)
    .then(result => {
      isSupportWebp = result;
    });
})();

export function supportWebp() {
  return isSupportWebp;
}

export function waitfor(ms) {
  const start = Date.now();
  return () => {
    const delay = ms - (Date.now() - start);
    if (delay < 0) {
      return;
    }
    return new Promise(resolve => setTimeout(resolve, delay));
  };
}

export function getCover(cover, height) {
  return `${coverUrlPrefix}/${cover}-90-0-${height}.jpeg`;
}

// getFontMetrics get font metrics instance
export function getFontMetrics(options) {
  const key = JSON.stringify(options);
  // 相关配置如果不变，font metrics可以重复使用，
  // 避免每次重新计算文本宽度
  if (currentFontMetrics && currentFontMetrics.key === key) {
    return currentFontMetrics.ins;
  }
  const fontMetrics = new FontMetrics(options);
  currentFontMetrics = {
    key,
    ins: fontMetrics
  };
  return fontMetrics;
}
