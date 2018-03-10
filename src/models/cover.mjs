import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const name = 'Cover';

const schema = {
  book: {
    type: Number,
    unique: true,
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
  client.model(name, s);
  return {
    name,
    schema: s,
  };
}
