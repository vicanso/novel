import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const schemaName = 'RequestBook';

const schema = {
  name: {
    type: String,
    required: true,
  },
  author: {
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
  client.model(schemaName, s);
  return {
    name: schemaName,
    schema: s,
  };
}
