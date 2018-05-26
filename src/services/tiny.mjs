import grpc from 'grpc';
import path from 'path';

import * as config from '../config';

const protoFile = path.join(config.appPath, '../protos/compress.proto');

const {compress} = grpc.load(protoFile);
const {WEBP, JPEG} = compress.Type;
const client = new compress.Compress(
  config.tiny,
  grpc.credentials.createInsecure(),
);

function convert(buf, type, quality) {
  const request = new compress.CompressRequest();
  request.setType(type);
  request.setData(new Uint8Array(buf));
  request.setQuality(quality);
  return new Promise((resolve, reject) => {
    client.do(request, (err, res) => {
      if (err) {
        reject(err);
      } else {
        resolve(res);
      }
    });
  });
}

/**
 * 转换为webp
 * @param buf
 */
export function toWebp(buf, quality) {
  return convert(buf, WEBP, quality);
}

/**
 * 转换为jpeg
 * @param buf
 */
export function toJpeg(buf, quality) {
  return convert(buf, JPEG, quality);
}
