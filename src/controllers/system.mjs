/**
 * 此模块主要是一些公共与业务无关的处理
 */

import moment from 'moment';
import util from 'util';
import path from 'path';
import fs from 'fs';
import _ from 'lodash';

import * as config from '../config';
import * as globals from '../helpers/globals';
import * as mongo from '../helpers/mongo';

const readFile = util.promisify(fs.readFile);

/**
 * 获取系统当前运行的版本package.json与读取文件package.json的版本号，
 * 正常情况下两者一致，但是如果更新了版本，但是没有重启就会不一致
 * @return {Object} 返回版本号信息 {pkg: 程序版本号, gen: 生成镜像的时间版本号}
 */
async function getVersion() {
  const buf = await readFile(path.join(config.appPath, '../assets/version'));
  return {
    pkg: config.version,
    gen: buf.toString(),
  };
}

/**
 * @swagger
 * /sys/pause:
 *  put:
 *    description: 暂停应用程序（只修改状态，程序并未停止）。中间件：m.admin
 *    summary: 暂停应用程序
 *    tags:
 *      - system
 *    responses:
 *      204:
 *        description: 执行成功时返回
 */
export function pause(ctx) {
  globals.pause();
  console.info('pause application');
  ctx.body = null;
}

/**
 * @swagger
 * /sys/resume:
 *  put:
 *    description: 恢复应用程序（将状态恢复）。中间件：m.admin
 *    summary: 恢复应用程序
 *    tags:
 *      - system
 *    responses:
 *      204:
 *        description: 执行成功时返回
 */
export function resume(ctx) {
  globals.start();
  console.info('resume application');
  ctx.body = null;
}

/**
 * @swagger
 * /sys/status:
 *  get:
 *    description: 获取应用信息，包括程序版本，连接数等。中间件：m.noQuery
 *    summary: 获取应用信息
 *    tags:
 *      - system
 *    responses:
 *      200:
 *        description: 应用程序状态信息
 *        schema:
 *          type: object
 *          properties:
 *            ins:
 *              type: string
 *              description: 随机生成的字符串，便于标记不同的应用实例
 *              example: HyOtXxZGf
 *            connectingTotal:
 *              description: 连接数
 *              type: integer
 *              example: 100
 *            status:
 *              description: 状态
 *              type: string
 *              example: running
 *              enum:
 *                - "running"
 *                - "pause"
 *            uptime:
 *              description: 运行时长
 *              type: string
 *              example: a minute ago
 *            startedAt:
 *              description: 启动时间
 *              type: string
 *              example: 2017-12-16T07:06:55.249Z
 *            pkg:
 *              description: 版本号
 *              type: string
 *              example: 1.0.0
 *            gen:
 *              description: 编译生成版本号（生成时间）
 *              type: string
 *              example: 2017-12-16T05:06:55.249Z
 */
export async function status(ctx) {
  const version = await getVersion();
  const sceonds = Number.parseInt(process.uptime(), 10);
  const uptime = moment(Date.now() - sceonds * 1000);
  ctx.setCache('10s');
  ctx.body = _.extend(
    {
      ins: config.ins,
      connectingTotal: globals.getConnectingCount(),
      status: globals.isRunning() ? 'running' : 'pause',
      uptime: uptime.fromNow(),
      startedAt: uptime.toISOString(),
    },
    version,
  );
}

/**
 * @swagger
 * /sys/stats:
 *  get:
 *    description: 获取应用程序性能统计，包括系统启动时间，使用CPU 内存等。中间件 m.noQuery
 *    summary: 获取应用程序性能统计
 *    tags:
 *      - system
 *    responses:
 *      200:
 *        description: 应用程序性能指标
 *        schema:
 *          type: object
 *          properties:
 *            app:
 *              description: 应用名称
 *              type: string
 *              example: spt
 *            uptime:
 *              description: 运行时间（秒）
 *              type: integer
 *              example: 134
 *            pid:
 *              description: 程序PID
 *              type: integer
 *              example: 9301
 *            connecting:
 *              description: 连接数
 *              type: integer
 *              example: 1,
 *            cpu:
 *              description: CPU使用率
 *              type: integer
 *              example: 10
 *            memory:
 *              description: 内存使用(MB)
 *              type: integer
 *              example: 52
 */
export async function stats(ctx) {
  const performance = globals.getPerformance();
  const data = {
    app: config.app,
    uptime: Math.ceil(process.uptime()),
    pid: process.pid,
  };
  if (performance) {
    Object.assign(data, {
      connecting: performance.connectingCount,
      cpu: performance.cpuUsageUsedPercent,
      memory: performance.memoryUsageRss,
    });
  }
  ctx.body = data;
}

/**
 * @swagger
 * /sys/exit:
 *  put:
 *    description: 退出应用程序，先设置状态为中止，再退出。中间件：m.admin
 *    summary: 退出应用程序
 *    tags:
 *      - system
 *    responses:
 *      204:
 *        description: 响应后直接返回，在10秒内完成退出
 *
 */
export function exit(ctx) {
  process.emit('SIGQUIT');
  ctx.body = null;
}

/**
 * 读取swagger api 文档信息
 */
export async function apis(ctx) {
  const file = path.join(config.appPath, 'assets/api.json');
  const buf = await readFile(file);
  ctx.setCache('5m');
  ctx.body = JSON.parse(buf);
}

/**
 * @swagger
 * /sys/:collection/ensure-indexes:
 *  put:
 *    description: 对mongodb中的collection确保索引创建。中间件：m.admin
 *    summary: 确保索引创建
 *    tags:
 *      - system
 *    responses:
 *      204:
 *        description: 执行成功时返回
 */
export async function ensureIndexes(ctx) {
  const {collection} = ctx.params;
  const Model = mongo.get(collection);
  await Model.ensureIndexes();
  ctx.body = null;
}

/**
 * @swagger
 * /sys/:collection/indexes:
 *  get:
 *    description: 获取该collection对应的索引
 *    summary: 获取索引
 *    tags:
 *      - system
 *    responses:
 *      200:
 *        description: 获取成功时返回
 */
export async function getIndexes(ctx) {
  const {collection} = ctx.params;
  const Model = mongo.get(collection);
  const data = await Model.collection.getIndexes();
  ctx.body = data;
}

// 返回管理页面
export async function adminIndex(ctx) {
  const file = path.join(config.appPath, '../admin/dist/index.html');
  const buf = await readFile(file);
  ctx.set('Content-Type', 'text/html; charset=utf-8');
  ctx.setCache('1m');
  // 根据系统启动环境替换前端配置
  const envScript = `var ENV = '${config.env}';`;
  ctx.body = buf.toString().replace("var ENV = 'development';", envScript);
}
