import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const name = 'Inc';

const schema = {
  category: {
    type: String,
    required: true,
    unique: true,
  },
  value: Number,
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
