const listBookPageSizeKey = "list-book-page-size";
export function saveListBookPageSize(size) {
  if (localStorage) {
    localStorage.setItem(listBookPageSizeKey, size);
  }
  return;
}

export function getListBookPageSize() {
  let pageSize = 10;
  if (localStorage) {
    const v = localStorage.getItem(listBookPageSizeKey);
    if (v) {
      pageSize = Number.parseInt(v, 10);
    }
  }
  return pageSize;
}
