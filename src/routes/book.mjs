export default [
  // 增加书籍来源
  [
    'POST',
    '/books/sources',
    ['m.admin', 'm.tracker("addSource")', 'c.book.addSource'],
  ],
  // 获取分类
  ['GET', '/books/categories', ['m.noQuery', 'c.book.categoriesList']],
  // 书籍列表
  ['GET', '/books', 'c.book.list'],
  // 增加书籍
  ['POST', '/books', ['m.tracker("addBook")', 'm.admin', 'c.book.add']],
  // 获取书籍信息
  ['GET', '/books/:no', 'c.book.get'],
  // 更新书籍
  [
    'PATCH',
    '/books/:no',
    ['m.admin', 'm.tracker("updateBook")', 'c.book.update'],
  ],
  // 更新书籍信息
  [
    'PATCH',
    '/books/:no/info',
    ['m.admin', 'm.tracker("updateBookInfo")', 'c.book.updateBookInfo'],
  ],
  // 获取章节
  ['GET', '/books/:no/chapters', 'c.book.listChapter'],
  // 获取封面
  ['GET', '/books/:no/cover', 'c.book.getCover'],
];
