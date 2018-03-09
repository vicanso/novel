import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const name = 'Origin';

const schema = {
  name: {
    type: String,
    required: true,
  },
  author: {
    type: String,
    required: true,
  },
  source: {
    type: String,
    required: true,
  },
  sourceId: {
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
