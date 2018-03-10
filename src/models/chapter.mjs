import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const name = 'Chapter';

const schema = {
  book: {
    type: Number,
    required: true,
  },
  no: {
    type: Number,
    required: true,
  },
  wordCount: {
    type: Number,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  data: {
    type: Buffer,
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
      book: 1,
    },
    {
      background: true,
    },
  );
  s.index(
    {
      book: 1,
      no: 1,
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
