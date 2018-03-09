import spider from 'novel-spider';

import genService from './gen';
import originService from './origin';
import {getBookNo} from './inc';

const bookService = genService('Book');

export default bookService;

export async function add(author, name) {
  let doc = await bookService.findOne({
    author,
    name,
  });
  if (doc) {
    return doc;
  }
  doc = await originService.findOne({
    name,
    author,
  });
  if (!doc || doc.source !== 'biquge') {
    return null;
  }
  const id = doc.sourceId;
  const novel = new spider.BiQuGe(id);
  const info = await novel.getInfos();
  if (!info || info.name !== name || info.author !== author) {
    return null;
  }
  const no = await getBookNo();
  const data = {
    no,
    name,
    author,
    brief: info.brief,
  };
  doc = await bookService.add(data);
  return doc;
}
