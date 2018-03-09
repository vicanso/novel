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
  latestChapter: {
    title: String,
    updatedAt: Date,
  },
  chapterCount: {
    type: Number,
    default: 0,
  },
  wordCount: {
    type: Number,
    default: 0,
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
  client.model(name, s);
  return {
    name,
    schema: s,
  };
}
