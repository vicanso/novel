import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const name = 'I18n';

const schema = {
  name: {
    type: String,
    required: true,
  },
  // 该字段分类，方便可以直接把所有分类的取出
  category: {
    type: String,
    required: true,
  },
  // 英文
  en: String,
  // 中文
  zh: String,
  creator: {
    type: String,
    required: true,
  },
};

export default function init(client) {
  const s = new Schema(schema, {
    timestamps: true,
    autoIndex: isDevelopment(),
  });
  s.index(
    {
      category: 1,
    },
    {
      background: true,
    },
  );
  s.index(
    {
      name: 1,
    },
    {
      background: true,
    },
  );
  s.index(
    {
      category: 1,
      name: 1,
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
