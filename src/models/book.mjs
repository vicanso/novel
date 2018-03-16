import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const name = 'Book';

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
};

export default function init(client) {
  const s = new Schema(schema, {
    timestamps: true,
    autoIndex: isDevelopment(),
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
  client.model(name, s);
  return {
    name,
    schema: s,
  };
}
