import mongoose from 'mongoose';

import {isDevelopment} from '../helpers/utils';

const {Schema} = mongoose;

const {Mixed} = Schema.Types;
const name = 'Setting';

const schema = {
  // 应用配置名称，用于区分不同的配置
  name: {
    type: String,
    unique: true,
  },
  // 该配置是否禁用状态
  disabled: {
    type: Boolean,
    default: false,
  },
  data: Mixed,
  creator: {
    type: String,
    required: true,
  },
  description: String,
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
