import mongoose from 'mongoose';

import {isDevelopment, getPinYin} from '../helpers/utils';

const {Schema} = mongoose;

const schemaName = 'Book';

const schema = {
  no: {
    type: Number,
    unique: true,
  },
  name: {
    type: String,
    required: true,
  },
  author: {
    type: String,
    required: true,
  },
  // 简介
  brief: String,
  // 类别
  category: [],
  // 最新章节
  latestChapter: {
    // 章节标题
    title: String,
    // 章节更新时间
    updatedAt: Date,
    // 章节字数
    wordCount: Number,
    // 章节序号
    no: Number,
  },
  chapterCount: {
    type: Number,
    default: 0,
  },
  wordCount: {
    type: Number,
    default: 0,
  },
  // 是否完结
  end: {
    type: Boolean,
    default: false,
  },
  // 关键字
  keyword: String,
};

export default function init(client) {
  const s = new Schema(schema, {
    timestamps: true,
    autoIndex: isDevelopment(),
  });
  s.pre('save', function preSave(next) {
    const {name, author} = this;
    if (name && author) {
      const keywords = [name, author];
      keywords.push(getPinYin(name));
      keywords.push(getPinYin(name, true));
      keywords.push(getPinYin(author));
      keywords.push(getPinYin(author, true));
      this.keyword = keywords.join(' ');
    }
    next();
  });
  s.index(
    {
      name: 1,
      author: 1,
    },
    {
      background: true,
      unique: true,
    },
  );
  s.index(
    {
      category: 1,
    },
    {
      background: true,
    },
  );
  client.model(schemaName, s);
  return {
    name: schemaName,
    schema: s,
  };
}
